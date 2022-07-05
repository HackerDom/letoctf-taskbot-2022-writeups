package adapter

import (
	"context"
	"sync/atomic"

	"github.com/jackc/pgx"

	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/basestore"
)

type PGXAdapter struct {
	pool       *pgx.ConnPool
	txnCounter uint64
}

func NewPGXAdapter(pool *pgx.ConnPool) *PGXAdapter {
	return &PGXAdapter{pool: pool}
}

func (d *PGXAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (basestore.ExecResult, error) {
	return d.pool.ExecEx(ctx, sql, nil, args...)
}

func (d *PGXAdapter) Query(ctx context.Context, sql string, args ...interface{}) (basestore.Rows, error) {
	return d.pool.QueryEx(ctx, sql, nil, args...)
}

func (d *PGXAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) basestore.Row {
	return d.pool.QueryRowEx(ctx, sql, nil, args...)
}

func (d *PGXAdapter) BeginTX(ctx context.Context, level basestore.IsolationLevel) (basestore.TXN, error) {
	var txLevel pgx.TxIsoLevel

	switch level {
	case basestore.ReadUncommitted:
		txLevel = pgx.ReadUncommitted
	case basestore.ReadCommitted:
		txLevel = pgx.ReadCommitted
	case basestore.RepeatableRead:
		txLevel = pgx.RepeatableRead
	case basestore.Serializable:
		txLevel = pgx.Serializable
	}

	tx, err := d.pool.BeginEx(
		ctx, &pgx.TxOptions{
			IsoLevel: txLevel,
		},
	)
	if err != nil {
		return nil, err
	}
	return &pgxTXN{
		tx:    tx,
		txnID: atomic.AddUint64(&d.txnCounter, 1),
	}, nil
}

type pgxTXN struct {
	tx    *pgx.Tx
	txnID uint64
}

func (txn *pgxTXN) Exec(ctx context.Context, sql string, args ...interface{}) (basestore.ExecResult, error) {
	return txn.tx.ExecEx(ctx, sql, nil, args...)
}

func (txn *pgxTXN) Query(ctx context.Context, sql string, args ...interface{}) (basestore.Rows, error) {
	return txn.tx.QueryEx(ctx, sql, nil, args...)
}

func (txn *pgxTXN) QueryRow(ctx context.Context, sql string, args ...interface{}) basestore.Row {
	return txn.tx.QueryRowEx(ctx, sql, nil, args...)
}

func (txn *pgxTXN) ID() basestore.TXID {
	return basestore.TXID(txn.txnID)
}

func (txn *pgxTXN) Commit(ctx context.Context) error {
	return txn.tx.CommitEx(ctx)
}

func (txn *pgxTXN) Rollback(ctx context.Context) error {
	return txn.tx.RollbackEx(ctx)
}
