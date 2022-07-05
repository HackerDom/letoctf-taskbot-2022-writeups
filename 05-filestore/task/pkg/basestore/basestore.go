package basestore

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type IsolationLevel int

const (
	ReadUncommitted IsolationLevel = iota
	ReadCommitted
	RepeatableRead
	Serializable
)

type Row interface {
	// Scan считывает значения колонок из базы данных.
	Scan(...interface{}) error
}

type Rows interface {
	// Scan считывает значения колонок из базы данных.
	Scan(...interface{}) error

	// Next подготавливает следующую строку для чтения методом Scan.
	// Возвращает false, если больше нет строк для считывания.
	Next() bool

	// Close завершает чтение колонок.
	Close()
}

type ExecResult interface {
	// RowsAffected возвращает количество строк затронутых SELECT, UPDATE, INSERT или DELETE выражениями.
	RowsAffected() int64
}

type Querier interface {
	// Exec выполняет запрос без возврата каких-либо строк.
	Exec(ctx context.Context, sql string, args ...interface{}) (ExecResult, error)

	// Query выполняет запрос, который возвращает строки, как правило, SELECT.
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)

	// QueryRow выполняет запрос, который должен вернуть не более одной строки. Всегда возвращает ненулевое значение.
	// Если в результате запроса не выбрано ни одной строки, то при вызове row.Scan будет возвращено значение ErrNoRows.
	// Если результат запрос несколько строк, то вычитывается первая, а остальные игнорируются.
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
}

type TXID uint64

type TXN interface {
	Querier

	ID() TXID

	// Commit применяет транзакцию.
	Commit(ctx context.Context) error

	// Rollback прерывает транзакцию.
	Rollback(ctx context.Context) error
}

type Adapter interface {
	Querier

	// BeginTX начинает транзакцию.
	BeginTX(ctx context.Context, level IsolationLevel) (TXN, error)
}

type BaseStore struct {
	adapter Adapter
}

func New(adapter Adapter) *BaseStore {
	return &BaseStore{
		adapter: adapter,
	}
}

func (s *BaseStore) RunInTransaction(ctx context.Context, f func(ctx context.Context) error) (err error) {
	return s.RunInTransactionWithLevel(ctx, ReadCommitted, f)
}

func (s *BaseStore) RunInTransactionWithLevel(
	ctx context.Context,
	level IsolationLevel,
	f func(ctx context.Context) error,
) (err error) {
	var finishTransactionFunc = func(_ error) error { return nil }
	logger := zerolog.Ctx(ctx).With().Str("subsys", "base_store").Logger()

	txn, ok := txnFromContext(ctx)
	if !ok {
		// Если мы не обнаруживаем транзакцию в контексте, то инициируем её.
		txn, err = s.adapter.BeginTX(ctx, level)
		if err != nil {
			return err
		}

		logger.Info().
			Str("tx_id", fmt.Sprint(txn.ID())).
			Msg("Create new transaction.")
		transactionStartTime := time.Now()

		// Добавим ссылку на транзакцию в контекст, чтобы следующие вызовы RunInTransaction его использовали.
		ctx = newTXNContext(ctx, txn)

		// Подготовим функцию, завершающую нашу транзакцию.
		finishTransactionFunc = func(err error) error {
			defer logger.Info().
				Str("tx_id", fmt.Sprint(txn.ID())).
				Str("elapsed_ms", fmt.Sprint(time.Since(transactionStartTime).Milliseconds())).
				Msg("Finish transaction.")
			if err != nil {
				return txn.Rollback(ctx)
			}

			return txn.Commit(ctx)
		}
	} else {
		logger.Info().
			Str("tx_id", fmt.Sprint(txn.ID())).
			Msg("Using existing transaction.")
	}

	if err = f(ctx); err != nil {
		if err := finishTransactionFunc(err); err != nil {
			// Напишем в лог, но вернём оригинальную ошибку пользователю.
			logger.Err(err).
				Str("tx_id", fmt.Sprint(txn.ID())).
				Msg("An error occurred at the time the transaction finished.")
		}
		return err
	}
	return finishTransactionFunc(nil)
}

func (s *BaseStore) Q(ctx context.Context) Querier {
	txn, ok := txnFromContext(ctx)
	if ok {
		return &txnQuerier{txn: txn}
	}
	return s.adapter
}

type txnQuerier struct {
	txn TXN
}

func (t *txnQuerier) Exec(ctx context.Context, sql string, args ...interface{}) (ExecResult, error) {
	return t.txn.Exec(ctx, sql, args...)
}

func (t *txnQuerier) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return t.txn.Query(ctx, sql, args...)
}

func (t *txnQuerier) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return t.txn.QueryRow(ctx, sql, args...)
}

type txContextKey struct{}

func newTXNContext(ctx context.Context, txn TXN) context.Context {
	return context.WithValue(ctx, txContextKey{}, txn)
}

func txnFromContext(ctx context.Context) (TXN, bool) {
	txn, ok := ctx.Value(txContextKey{}).(TXN)
	return txn, ok
}
