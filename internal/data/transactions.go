package data

import "gitlab.com/distributed_lab/kit/pgdb"

type TransactionsQ interface {
	New() TransactionsQ

	Get() (*Transaction, error)
	Select() ([]Transaction, error)
	Exists(requestID int64, status TxStatus) (bool, error)
	Update() ([]Transaction, error)

	Transaction(fn func(q TransactionsQ) error) error

	Page(pageParams pgdb.OffsetPageParams) TransactionsQ

	FilterByRequestID(ids ...int64) TransactionsQ
	FilterByHash(hashes ...string) TransactionsQ
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

type Transaction struct {
	Hash          string   `db:"hash" structs:"hash"`
	Body          string   `db:"body" structs:"body"`
	RequestID     int64    `db:"request_id" structs:"request_id"`
	Status        TxStatus `db:"status" structs:"status"`
	FailureReason *string  `db:"failure_reason" structs:"failure_reason"`
}
