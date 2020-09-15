package pg

import (
	"database/sql"
	"errors"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const requestsTableName = "requests"

func NewRequestsQ(db *pgdb.DB) data.RequestsQ {
	return &requestsQ{
		db:        db.Clone(),
		sql:       sq.Select("*").From(requestsTableName),
		sqlUpdate: sq.Update(requestsTableName),
	}
}

type requestsQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (q *requestsQ) New() data.RequestsQ {
	return NewRequestsQ(q.db)
}

func (q *requestsQ) Get() (*data.Request, error) {
	var result data.Request
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *requestsQ) Select() ([]data.Request, error) {
	var result []data.Request
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *requestsQ) Update() ([]data.Request, error) {
	var result []data.Request
	err := q.db.Select(&result, q.sqlUpdate)
	return result, err
}

func (q *requestsQ) Transaction(fn func(q data.RequestsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *requestsQ) Insert(requests ...data.Request) ([]data.Request, error) {
	if len(requests) == 0 {
		return nil, errors.New("empty array is not allowed")
	}

	names := []string{
		"id",
		"owner",
		"status",
		"failure_reason",
	}
	stmt := sq.Insert(requestsTableName).Columns(names...)
	for _, item := range requests {
		stmt = stmt.Values([]interface{}{
			item.ID,
			item.Owner,
			item.Status,
			item.FailureReason,
		}...)
	}

	stmt = stmt.Suffix("returning *")
	var result []data.Request
	err := q.db.Select(&result, stmt)

	return result, err
}

func (q *requestsQ) Page(pageParams pgdb.OffsetPageParams) data.RequestsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *requestsQ) FilterByID(ids ...int64) data.RequestsQ {
	q.sql = q.sql.Where(sq.Eq{"id": ids})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"id": ids})
	return q
}

func (q *requestsQ) FilterByOwner(owners ...string) data.RequestsQ {
	q.sql = q.sql.Where(sq.Eq{"owner": owners})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"owner": owners})
	return q
}

func (q *requestsQ) FilterByStatus(statuses ...data.RequestStatus) data.RequestsQ {
	q.sql = q.sql.Where(sq.Eq{"status": statuses})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"status": statuses})
	return q
}

func (q *requestsQ) SetStatus(status data.RequestStatus) data.RequestsQ {
	q.sqlUpdate = q.sqlUpdate.Set("status", status)
	return q
}

func (q *requestsQ) SetFailureReason(reason string) data.RequestsQ {
	q.sqlUpdate = q.sqlUpdate.Set("failure_reason", reason)
	return q
}

func (q *requestsQ) GetMaxId() (*int64, error) {
	stmt := sq.Select("max(id)").From(requestsTableName)

	var result *int64
	err := q.db.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func (q *requestsQ) InsertPayments(txs ...data.Transaction) ([]data.Transaction, error) {
	if len(txs) == 0 {
		return nil, errors.New("empty array is not allowed")
	}

	names := []string{
		"request_id",
		"status",
		"failure_reason",
		"amount",
		"destination",
		"destination_type",
	}
	stmt := sq.Insert(transactionsTableName).Columns(names...)
	for _, item := range txs {
		stmt = stmt.Values([]interface{}{
			item.RequestID,
			item.Status,
			item.FailureReason,
			item.Amount,
			item.Destination,
		}...)
	}

	stmt = stmt.Suffix("returning *")
	var result []data.Transaction
	err := q.db.Select(&result, stmt)

	return result, err
}
