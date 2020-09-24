package streamers

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/go/xdrbuild"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/connectors/request"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/horizon"
	regources "gitlab.com/tokend/regources/generated"
)

func NewCreateDeferredPaymentRequestProcessor(log *logan.Entry, requestsQ data.RequestsQ, reviewer *request.Reviewer,
	client *horizon.Connector, tasks uint32) request.PageProcessor {
	return &createDeferredPaymentRequestProcessor{
		log:                    log,
		requestsQ:              requestsQ,
		reviewer:               reviewer,
		client:                 client,
		reviewableRequestTasks: tasks,
	}
}

type createDeferredPaymentRequestProcessor struct {
	log                    *logan.Entry
	requestsQ              data.RequestsQ
	reviewer               *request.Reviewer
	client                 *horizon.Connector
	reviewableRequestTasks uint32
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
	_, err := p.reviewer.Reject(context.TODO(), r,
		request.ReviewDetails{RejectReason: reason, ExternalDetails: "{}"},
		xdrbuild.ReviewableRequestBaseDetails{RequestType: xdr.ReviewableRequestTypeCreateDeferredPayment})
	if err != nil {
		return errors.Wrap(err, "failed to reject request")
	}

	return nil
}

func (p *createDeferredPaymentRequestProcessor) approve(r regources.ReviewableRequest) (*xdr.ReviewRequestResult, error) {
	return p.reviewer.Approve(context.TODO(), r, request.ReviewDetails{
		TasksToRemove:   p.reviewableRequestTasks,
		ExternalDetails: "{}",
	}, xdrbuild.ReviewableRequestBaseDetails{RequestType: xdr.ReviewableRequestTypeCreateDeferredPayment})
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

	res, err := p.approve(r)
	deferredPaymentID := int64(res.Success.TypeExt.CreateDeferredPaymentResult.DeferredPaymentId)

	paymentsBatches := convertToDbPayments(deferredPaymentID, rawPaymentsBatches)

	return p.saveRequest(data.Request{
		ID:          deferredPaymentID,
		Owner:       r.Relationships.Requestor.Data.ID,
		Status:      data.RequestStatusProcessing,
		LockupUntil: content.LockupUntil,
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
		for i, payment := range paymentsBatch {
			resultBatch[i] = data.Payment{
				RequestID:       requestId,
				Status:          data.PaymentStatusProcessing,
				Amount:          payment.Amount,
				Destination:     payment.Destination,
				DestinationType: payment.DestinationType,
			}
		}
		result = append(result, resultBatch)
	}
	return nil
}

type detailsContent struct {
	Blobs       []string   `json:"blobs"`
	LockupUntil *time.Time `json:"lockup_until"`
}

type rawPayment struct {
	Amount          regources.Amount `json:"amount"`
	Destination     string           `json:"destination"`
	DestinationType string           `json:"destination_type"`
}
