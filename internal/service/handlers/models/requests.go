package models

import (
	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"
	"gitlab.com/tokend/mass-payments-sender-svc/resources"
)

func NewRequestResponse(request data.Request) resources.RequestResponse {
	return resources.RequestResponse{
		Data: newRequestResponseModel(request),
	}
}

func NewRequestListResponse(requests []data.Request) resources.RequestListResponse {
	result := resources.RequestListResponse{
		Data: make([]resources.Request, len(requests)),
	}
	for i, request := range requests {
		result.Data[i] = newRequestResponseModel(request)
	}

	return result
}

func newRequestResponseModel(request data.Request) resources.Request {
	return resources.Request{
		Key: resources.NewKeyInt64(request.ID, resources.MASS_PAYMENTS_REQUESTS),
		Relationships: resources.RequestRelationships{
			Owner: *resources.Key{
				ID:   request.Owner,
				Type: resources.ACCOUNTS,
			}.AsRelation(),
		},
		Attributes: resources.RequestAttributes{
			Status:        string(request.Status),
			FailureReason: request.FailureReason,
			LockupUntil:   request.LockupUntil,
		},
	}
}
