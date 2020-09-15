package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type RequestsQ interface {
	New() RequestsQ

	Get() (*Request, error)
	Select() ([]Request, error)
	Insert(requests ...Request) ([]Request, error)
	Update() ([]Request, error)

	Transaction(fn func(q RequestsQ) error) error

	Page(pageParams pgdb.OffsetPageParams) RequestsQ

	FilterByID(ids ...int64) RequestsQ
	FilterByOwner(owners ...string) RequestsQ
	FilterByStatus(statuses ...RequestStatus) RequestsQ

	SetStatus(status RequestStatus) RequestsQ
	SetFailureReason(reason string) RequestsQ

	GetMaxId() (*int64, error)

	InsertPayments(payments ...Transaction) ([]Transaction, error)
}

type RequestStatus string

const (
	RequestStatusProcessing RequestStatus = "processing"
	RequestStatusFailed     RequestStatus = "failed"
	RequestStatusFinished   RequestStatus = "finished"
)

type Request struct {
	ID            int64         `db:"id" structs:"id"`
	Owner         string        `db:"owner" structs:"owner"`
	Status        RequestStatus `db:"status" structs:"status"`
	FailureReason *string       `db:"failure_reason" structs:"failure_reason"`
	LockupUntil   *time.Time    `db:"lockup_until" structs:"lockup_until"`
}
