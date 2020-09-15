package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	regources "gitlab.com/tokend/regources/generated"
)

type TransactionsQ interface {
	New() TransactionsQ

	Get() (*Transaction, error)
	Select() ([]Transaction, error)
	Exists(requestID int64, status TxStatus) (bool, error)
	Update() ([]Transaction, error)

	Transaction(fn func(q TransactionsQ) error) error

	Page(pageParams pgdb.OffsetPageParams) TransactionsQ

	FilterByRequestID(ids ...int64) TransactionsQ
	FilterByID(ids ...int64) TransactionsQ
	FilterByStatus(statuses ...TxStatus) TransactionsQ
	Limit(limit uint64) TransactionsQ

	SetStatus(status TxStatus) TransactionsQ
	SetFailureReason(reason string) TransactionsQ

	UpdateRequestStatus(requestId int64, status RequestStatus) error
}

type TxStatus string

const (
	TxStatusProcessing TxStatus = "processing"
	TxStatusFailed     TxStatus = "failed"
	TxStatusSuccess    TxStatus = "success"
)

const (
	DestinationTypeAccountID = "account_id"
)

type Transaction struct {
	ID              string           `db:"id" structs:"-"`
	RequestID       int64            `db:"request_id" structs:"request_id"`
	Status          TxStatus         `db:"status" structs:"status"`
	FailureReason   *string          `db:"failure_reason" structs:"failure_reason"`
	Amount          regources.Amount `db:"amount" structs:"amount"`
	Destination     string           `db:"destination" structs:"destination"`
	DestinationType string           `db:"destination_type" structs:"destination_type"`
}
