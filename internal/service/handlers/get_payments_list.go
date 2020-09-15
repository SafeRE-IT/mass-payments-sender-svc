package handlers

import (
	"net/http"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/handlers/models"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetPaymentsList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetPaymentsListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	dataOwners := make([]string, 0)
	if request.FilterRequestId != nil {
		req, err := RequestsQ(r).FilterByID(*request.FilterRequestId).Get()
		if err != nil {
			Log(r).WithError(err).Error("failed to get transactions from db")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if req != nil {
			dataOwners = append(dataOwners, req.Owner)
		}
	}

	if !IsAllowed(r, w, dataOwners...) {
		return
	}

	q := TransactionsQ(r).Page(request.OffsetPageParams)

	if request.FilterRequestId != nil {
		q.FilterByRequestID(*request.FilterRequestId)
	}

	if len(request.FilterStatus) > 0 {
		statuses := make([]data.PaymentStatus, len(request.FilterStatus))
		for i, status := range request.FilterStatus {
			statuses[i] = data.PaymentStatus(status)
		}
		q.FilterByStatus(statuses...)
	}

	txs, err := q.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get transactions from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.NewPaymentsResponseList(txs)
	response.Links = GetOffsetLinks(r, request.OffsetPageParams)
	ape.Render(w, response)
}
