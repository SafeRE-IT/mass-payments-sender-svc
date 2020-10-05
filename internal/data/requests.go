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
	FilterBySourceBalance(sourceBalances ...string) RequestsQ
	FilterByAsset(assets ...string) RequestsQ
	FilterByCreatedAt(from *time.Time, to *time.Time) RequestsQ

	SetStatus(status RequestStatus) RequestsQ
	SetFailureReason(reason string) RequestsQ

	GetMaxId() (*int64, error)

	InsertPayments(payments ...Payment) ([]Payment, error)
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
	SourceBalance string        `db:"source_balance" structs:"source_balance"`
	Asset         string        `db:"asset" structs:"asset"`
	Status        RequestStatus `db:"status" structs:"status"`
	FailureReason *string       `db:"failure_reason" structs:"failure_reason"`
	CreatedAt     time.Time     `db:"created_at" structs:"created_at"`
	LockupUntil   *time.Time    `db:"lockup_until" structs:"lockup_until"`
}
