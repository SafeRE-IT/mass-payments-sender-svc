package request

import (
"context"
"github.com/spf13/cast"
"gitlab.com/distributed_lab/logan/v3/errors"
"gitlab.com/tokend/connectors/keyer"
"gitlab.com/tokend/connectors/submit"
"gitlab.com/tokend/go/xdr"
"gitlab.com/tokend/go/xdrbuild"
regources "gitlab.com/tokend/regources/generated"
)

type ReviewDetails struct {
	TasksToAdd      uint32
	TasksToRemove   uint32
	ExternalDetails string

	RejectReason string
}

type ReviewerConfig interface {
	Builder() *xdrbuild.Builder
	Keys() keyer.Keys
	Submit() *submit.Submitter
}

type Reviewer struct {
	keys      keyer.Keys
	submitter *submit.Submitter
	builder   *xdrbuild.Builder
}

func NewReviewer(config ReviewerConfig) *Reviewer {
	return &Reviewer{
		keys:      config.Keys(),
		submitter: config.Submit(),
		builder:   config.Builder(),
	}
}

func (s *Reviewer) Approve(ctx context.Context,
	request regources.ReviewableRequest,
	reviewDetails ReviewDetails,
	externalDetails xdrbuild.ReviewRequestDetailsProvider) (*xdr.ReviewRequestResult, error) {
	id := cast.ToUint64(request.ID)

	tx := s.builder.Transaction(s.keys.Source).Op(&xdrbuild.ReviewRequest{
		ID:      id,
		Hash:    &request.Attributes.Hash,
		Action:  xdr.ReviewRequestOpActionApprove,
		Reason:  reviewDetails.RejectReason,
		Details: externalDetails,
		ReviewDetails: xdrbuild.ReviewDetails{
			TasksToAdd:      reviewDetails.TasksToAdd,
			TasksToRemove:   reviewDetails.TasksToRemove,
			ExternalDetails: reviewDetails.ExternalDetails,
		},
	}).Sign(s.keys.Signer)

	return s.submitTransaction(ctx, tx)
}

func (s *Reviewer) Reject(ctx context.Context,
	request regources.ReviewableRequest,
	reviewDetails ReviewDetails,
	externalDetails xdrbuild.ReviewRequestDetailsProvider) (*xdr.ReviewRequestResult, error) {
	id := cast.ToUint64(request.ID)

	tx := s.builder.Transaction(s.keys.Source).Op(&xdrbuild.ReviewRequest{
		ID:      id,
		Hash:    &request.Attributes.Hash,
		Action:  xdr.ReviewRequestOpActionReject,
		Reason:  reviewDetails.RejectReason,
		Details: externalDetails,
		ReviewDetails: xdrbuild.ReviewDetails{
			TasksToAdd:      reviewDetails.TasksToAdd,
			TasksToRemove:   reviewDetails.TasksToRemove,
			ExternalDetails: reviewDetails.ExternalDetails,
		},
	}).Sign(s.keys.Signer)

	return s.submitTransaction(ctx, tx)
}

func (s *Reviewer) PermanentReject(ctx context.Context,
	request regources.ReviewableRequest,
	reviewDetails ReviewDetails,
	externalDetails xdrbuild.ReviewRequestDetailsProvider) (*xdr.ReviewRequestResult, error) {
	id := cast.ToUint64(request.ID)

	tx := s.builder.Transaction(s.keys.Source).Op(&xdrbuild.ReviewRequest{
		ID:      id,
		Hash:    &request.Attributes.Hash,
		Action:  xdr.ReviewRequestOpActionPermanentReject,
		Reason:  reviewDetails.RejectReason,
		Details: externalDetails,
		ReviewDetails: xdrbuild.ReviewDetails{
			TasksToAdd:      reviewDetails.TasksToAdd,
			TasksToRemove:   reviewDetails.TasksToRemove,
			ExternalDetails: reviewDetails.ExternalDetails,
		},
	}).Sign(s.keys.Signer)

	return s.submitTransaction(ctx, tx)
}

func (s *Reviewer) submitTransaction(ctx context.Context, tx *xdrbuild.Transaction) (*xdr.ReviewRequestResult, error) {
	env, err := tx.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal transaction")
	}

	res, err := s.submitter.Submit(ctx, env, false)
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
