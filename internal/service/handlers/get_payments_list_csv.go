package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/requests"
)

func GetPaymentsListCsv(w http.ResponseWriter, r *http.Request) {
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

	q := TransactionsQ(r)

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

	renderPaymentsCsv(w, txs)
}

func renderPaymentsCsv(w http.ResponseWriter, txs []data.Payment) {
	records := make([][]string, 0, len(txs)+1)
	headers := []string{
		"ID",
		"Request ID",
		"Amount",
		"Destination",
		"Destination type",
		"Status",
		"Failure reason",
	}
	records = append(records, headers)
	for _, payment := range txs {
		raw := []string{
			fmt.Sprintf("%d", payment.ID),
			fmt.Sprintf("%d", payment.RequestID),
			payment.Amount.String(),
			payment.Destination,
			payment.DestinationType,
			string(payment.Status),
			toStringSafe(payment.FailureReason),
		}

		records = append(records, raw)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", "payments.csv"))
	csv.NewWriter(w).WriteAll(records)
}

// toStringSafe - cast pointer string to string in case of null returns string empty value
func toStringSafe(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
