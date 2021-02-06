package streamers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"

	"gitlab.com/tokend/connectors/submit"

	"github.com/spf13/cast"

	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/cosigner"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/go/xdrbuild"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/connectors/request"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/horizon"
	regources "gitlab.com/tokend/regources/generated"
)

func NewCreateDeferredPaymentRequestProcessor(log *logan.Entry, requestsQ data.RequestsQ,
	client *horizon.Connector, tasks uint32, cosigner cosigner.Cosigner, connector *horizon.Connector,
	builder *xdrbuild.Builder, signer keypair.Full, source keypair.Address) request.PageProcessor {
	return &createDeferredPaymentRequestProcessor{
		log:                    log,
		requestsQ:              requestsQ,
		client:                 client,
		reviewableRequestTasks: tasks,
		cosigner:               cosigner,
		connector:              connector,
		builder:                builder,
		signer:                 signer,
		source:                 source,
	}
}

type createDeferredPaymentRequestProcessor struct {
	log                    *logan.Entry
	requestsQ              data.RequestsQ
	client                 *horizon.Connector
	reviewableRequestTasks uint32
	cosigner               cosigner.Cosigner
	connector              *horizon.Connector
	builder                *xdrbuild.Builder
	signer                 keypair.Full
	source                 keypair.Address
}

func (p *createDeferredPaymentRequestProcessor) ProcessPage(ctx context.Context, page regources.ReviewableRequestListResponse) error {
	for _, r := range page.Data {
		p.log.Debugf("processing request with id: %s", r.ID)
		details := page.Included.MustCreateDeferredPaymentRequest(r.Relationships.RequestDetails.Data.GetKey())
		if details == nil {
			return errors.New("details not found")
		}
		err := p.processDeferredPayment(r, *details)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *createDeferredPaymentRequestProcessor) reject(r regources.ReviewableRequest, reason string) error {
	_, err := p.review(context.TODO(), r,
		request.ReviewDetails{RejectReason: reason, ExternalDetails: "{}"},
		xdrbuild.ReviewableRequestBaseDetails{RequestType: xdr.ReviewableRequestTypeCreateDeferredPayment},
		xdr.ReviewRequestOpActionReject)
	if err != nil {
		return errors.Wrap(err, "failed to reject request")
	}

	return nil
}

func (p *createDeferredPaymentRequestProcessor) approve(r regources.ReviewableRequest) (*xdr.ReviewRequestResult, error) {
	return p.review(context.TODO(), r, request.ReviewDetails{
		TasksToRemove:   p.reviewableRequestTasks,
		ExternalDetails: "{}",
	}, xdrbuild.ReviewableRequestBaseDetails{RequestType: xdr.ReviewableRequestTypeCreateDeferredPayment},
		xdr.ReviewRequestOpActionApprove)
}

func (p *createDeferredPaymentRequestProcessor) review(ctx context.Context,
	request regources.ReviewableRequest,
	reviewDetails request.ReviewDetails,
	externalDetails xdrbuild.ReviewRequestDetailsProvider,
	action xdr.ReviewRequestOpAction) (*xdr.ReviewRequestResult, error) {

	id := cast.ToUint64(request.ID)

	tx := p.builder.Transaction(p.source).Op(&xdrbuild.ReviewRequest{
		ID:      id,
		Hash:    &request.Attributes.Hash,
		Action:  action,
		Reason:  reviewDetails.RejectReason,
		Details: externalDetails,
		ReviewDetails: xdrbuild.ReviewDetails{
			TasksToAdd:      reviewDetails.TasksToAdd,
			TasksToRemove:   reviewDetails.TasksToRemove,
			ExternalDetails: reviewDetails.ExternalDetails,
		},
	}).Sign(p.signer)
	txBase64, err := tx.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal tx to base64")
	}

	txBase64, err = p.cosigner.Cosign(txBase64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cosign tx")
	}

	res, err := p.connector.Submit(ctx, txBase64, false)
	if err != nil {
		if serr, ok := err.(submit.TxFailure); ok {
			err = errors.From(err, serr.GetLoganFields())
		}
		return nil, errors.Wrap(err, "failed to submit transaction")
	}

	var txResult xdr.TransactionResult
	err = xdr.SafeUnmarshalBase64(res.Data.Attributes.ResultXdr, &txResult)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal transaction result")
	}

	return (*txResult.Result.Results)[0].Tr.ReviewRequestResult, nil
}

func (p *createDeferredPaymentRequestProcessor) processDeferredPayment(r regources.ReviewableRequest,
	details regources.CreateDeferredPaymentRequest) error {

	if details.Attributes.CreatorDetails == nil {
		return p.reject(r, "wrong data content")
	}

	var content detailsContent
	err := json.Unmarshal(*details.Attributes.CreatorDetails, &content)
	if err != nil {
		return p.reject(r, "wrong data content")
	}

	rawPaymentsBatches, err := p.downloadPayments(content, r)
	if err != nil {
		return errors.Wrap(err, "failed to download blob")
	}
	if rawPaymentsBatches == nil {
		// Request was rejected
		return nil
	}

	sourceBalance, err := p.client.GetBalance(details.Relationships.SourceBalance.Data.ID)
	if err != nil {
		return errors.Wrap(err, "failed to download source balance")
	}
	if sourceBalance == nil {
		return errors.New("source balance not found")
	}

	res, err := p.approve(r)
	if err != nil {
		return errors.Wrap(err, "failed to approve request")
	}
	deferredPaymentID := int64(res.Success.TypeExt.CreateDeferredPaymentResult.DeferredPaymentId)

	paymentsBatches := convertToDbPayments(deferredPaymentID, rawPaymentsBatches)

	return p.saveRequest(data.Request{
		ID:            deferredPaymentID,
		Owner:         r.Relationships.Requestor.Data.ID,
		Status:        data.RequestStatusProcessing,
		SourceBalance: sourceBalance.Data.ID,
		Asset:         sourceBalance.Data.Relationships.Asset.Data.ID,
		LockupUntil:   content.LockupUntil,
		CreatedAt:     time.Now().UTC(),
	}, paymentsBatches)
}

func (p *createDeferredPaymentRequestProcessor) downloadPayments(content detailsContent, r regources.ReviewableRequest) ([][]rawPayment, error) {
	result := make([][]rawPayment, 0, len(content.Blobs))
	for _, id := range content.Blobs {
		blob, err := p.client.GetBlob(id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to download blob")
		}

		var rawPayments []rawPayment
		err = json.Unmarshal([]byte(blob.Data.Attributes.Value), &rawPayments)
		if err != nil {
			return nil, p.reject(r, "wrong content")
		}

		result = append(result, rawPayments)
	}

	return result, nil
}

func (p *createDeferredPaymentRequestProcessor) saveRequest(request data.Request, paymentsBatches [][]data.Payment) error {
	return p.requestsQ.Transaction(func(q data.RequestsQ) error {
		_, err := q.Insert(request)
		if err != nil {
			return errors.Wrap(err, "failed to save request")
		}

		for _, batch := range paymentsBatches {
			_, err := q.InsertPayments(batch...)
			if err != nil {
				return errors.Wrap(err, "failed to save transactions batch")
			}
		}

		return nil
	})
}

func convertToDbPayments(requestId int64, rawPayments [][]rawPayment) [][]data.Payment {
	result := make([][]data.Payment, 0)
	for _, paymentsBatch := range rawPayments {
		resultBatch := make([]data.Payment, len(paymentsBatch))
		for i, rawPayment := range paymentsBatch {
			payment := data.Payment{
				RequestID:       requestId,
				Status:          data.PaymentStatusProcessing,
				Amount:          rawPayment.Amount,
				Destination:     rawPayment.Destination,
				DestinationType: rawPayment.DestinationType,
			}

			if rawPayment.CreatorDetails != nil {
				payment.CreatorDetails = types.NullJSONText{
					Valid:    true,
					JSONText: types.JSONText(*rawPayment.CreatorDetails),
				}
			}

			resultBatch[i] = payment
		}
		result = append(result, resultBatch)
	}
	return result
}

type detailsContent struct {
	Blobs       []string   `json:"blobs"`
	LockupUntil *time.Time `json:"lockup_until"`
}

type rawPayment struct {
	Amount          regources.Amount `json:"amount"`
	Destination     string           `json:"destination"`
	DestinationType string           `json:"destination_type"`
	CreatorDetails  *json.RawMessage `json:"creator_details,omitempty"`
}
