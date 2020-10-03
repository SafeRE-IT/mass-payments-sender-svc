package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	regources "gitlab.com/tokend/regources/generated"
)

type PaymentsQ interface {
	New() PaymentsQ

	Get() (*Payment, error)
	Select() ([]Payment, error)
	Exists(requestID int64, status PaymentStatus) (bool, error)
	Update() ([]Payment, error)

	Transaction(fn func(q PaymentsQ) error) error

	Page(pageParams pgdb.OffsetPageParams) PaymentsQ

	FilterByRequestID(ids ...int64) PaymentsQ
	FilterByID(ids ...int64) PaymentsQ
	FilterByStatus(statuses ...PaymentStatus) PaymentsQ
	Limit(limit uint64) PaymentsQ

	SetStatus(status PaymentStatus) PaymentsQ
	SetFailureReason(reason string) PaymentsQ

	UpdateRequestStatus(requestId int64, status RequestStatus) error
}

type PaymentStatus string

const (
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusSuccess    PaymentStatus = "success"
)

const (
	DestinationTypeAccountID = "account_id"
)

type Payment struct {
	ID              int64            `db:"id" structs:"-"`
	RequestID       int64            `db:"request_id" structs:"request_id"`
	Status          PaymentStatus    `db:"status" structs:"status"`
	FailureReason   *string          `db:"failure_reason" structs:"failure_reason"`
	Amount          regources.Amount `db:"amount" structs:"amount"`
	Destination     string           `db:"destination" structs:"destination"`
	DestinationType string           `db:"destination_type" structs:"destination_type"`
	TxBody          *string          `db:"tx_body" structs:"tx_body"`
}
