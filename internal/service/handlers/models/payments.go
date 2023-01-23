package models

import (
	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/data"
	"github.com/SafeRE-IT/mass-payments-sender-svc/resources"
)

func NewPaymentsResponseList(models []data.Payment) resources.PaymentListResponse {
	result := resources.PaymentListResponse{
		Data: make([]resources.Payment, len(models)),
	}
	for i, model := range models {
		result.Data[i] = newPaymentsResponseModel(model)
	}

	return result
}

func NewPaymentsResponse(model data.Payment) resources.PaymentResponse {
	return resources.PaymentResponse{
		Data: newPaymentsResponseModel(model),
	}
}

func newPaymentsResponseModel(model data.Payment) resources.Payment {
	return resources.Payment{
		Key: resources.NewKeyInt64(model.ID, resources.MASS_PAYMENTS_PAYMENTS),
		Attributes: resources.PaymentAttributes{
			Status:          string(model.Status),
			FailureReason:   model.FailureReason,
			Amount:          model.Amount,
			Destination:     model.Destination,
			DestinationType: model.DestinationType,
		},
		Relationships: resources.PaymentRelationships{
			Request: *resources.NewKeyInt64(model.RequestID, resources.MASS_PAYMENTS_REQUESTS).AsRelation(),
		},
	}
}
