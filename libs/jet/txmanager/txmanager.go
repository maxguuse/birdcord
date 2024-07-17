package txmanager

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrTx   = errors.New("tx error")
	ErrFunc = errors.New("func error")
)

type ctxKey struct{}

var defaultCtxKey = ctxKey{}

type TxManager struct {
	defaultDb *sql.DB
}

func NewManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{
		defaultDb: stdlib.OpenDBFromPool(db),
	}
}

func (txm *TxManager) Do(
	ctx context.Context,
	f func(context.Context) error,
) error {
	tx, err := txm.defaultDb.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(
			errors.New("txm.defaultDb.Begin: "),
			ErrTx,
			err,
		)
	}

	defer tx.Rollback() //nolint: errcheck

	err = f(context.WithValue(ctx, defaultCtxKey, tx))

	if err != nil {
		return errors.Join(
			errors.New("f: "),
			ErrFunc,
			err,
		)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(
			errors.New("tx.Commit: "),
			ErrTx,
			err,
		)
	}

	return nil
}

type TxGetter struct {
	txManager *TxManager
}

func NewGetter(txm *TxManager) *TxGetter {
	return &TxGetter{
		txManager: txm,
	}
}

func (txg *TxGetter) DefaultTxOrDB(ctx context.Context) qrm.DB { //nolint: ireturn
	val := ctx.Value(defaultCtxKey)
	if val == nil {
		return txg.txManager.defaultDb
	}

	tx, valid := val.(*sql.Tx)
	if !valid {
		return txg.txManager.defaultDb
	}

	return tx
}
