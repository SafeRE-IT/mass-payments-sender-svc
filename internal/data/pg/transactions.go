package pg

import (
	"database/sql"
	"fmt"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const transactionsTableName = "transactions"

func NewTransactionsQ(db *pgdb.DB) data.TransactionsQ {
	return &transactionsQ{
		db:        db.Clone(),
		sql:       sq.Select("*").From(transactionsTableName),
		sqlUpdate: sq.Update(transactionsTableName),
	}
}

type transactionsQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (q *transactionsQ) New() data.TransactionsQ {
	return NewTransactionsQ(q.db)
}

func (q *transactionsQ) Get() (*data.Transaction, error) {
	var result data.Transaction
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *transactionsQ) Select() ([]data.Transaction, error) {
	var result []data.Transaction
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *transactionsQ) Exists(requestID int64, status data.TxStatus) (bool, error) {
	stmt := sq.Select(fmt.Sprintf("exists(select 1 from %s where request_id = %d and status = '%s')",
		transactionsTableName, requestID, status))

	var result bool
	err := q.db.Get(&result, stmt)
	return result, err
}

func (q *transactionsQ) Update() ([]data.Transaction, error) {
	var result []data.Transaction
	err := q.db.Select(&result, q.sqlUpdate)
	return result, err
}

func (q *transactionsQ) Transaction(fn func(q data.TransactionsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *transactionsQ) Page(pageParams pgdb.OffsetPageParams) data.TransactionsQ {
	q.sql = pageParams.ApplyTo(q.sql, "hash")
	return q
}

func (q *transactionsQ) FilterByRequestID(ids ...int64) data.TransactionsQ {
	q.sql = q.sql.Where(sq.Eq{"request_id": ids})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"request_id": ids})
	return q
}

func (q *transactionsQ) FilterByHash(hashes ...string) data.TransactionsQ {
	q.sql = q.sql.Where(sq.Eq{"hash": hashes})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"hash": hashes})
	return q
}

func (q *transactionsQ) FilterByStatus(statuses ...data.TxStatus) data.TransactionsQ {
	q.sql = q.sql.Where(sq.Eq{"status": statuses})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"status": statuses})
	return q
}

func (q *transactionsQ) Limit(limit uint64) data.TransactionsQ {
	q.sql = q.sql.Limit(limit)
	q.sqlUpdate = q.sqlUpdate.Limit(limit)
	return q
}

func (q *transactionsQ) SetStatus(status data.TxStatus) data.TransactionsQ {
	q.sqlUpdate = q.sqlUpdate.Set("status", status)
	return q
}

func (q *transactionsQ) SetFailureReason(reason string) data.TransactionsQ {
	q.sqlUpdate = q.sqlUpdate.Set("failure_reason", reason)
	return q
}

func (q *transactionsQ) UpdateRequestStatus(requestId int64, status data.RequestStatus) error {
	stmt := sq.Update(requestsTableName).
		Set("status", status).
		Where(sq.Eq{"id": requestId})

	return q.db.Exec(stmt)
}
