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

type TxManager struct {
	defaultDb *sql.DB
}

func New(db *pgxpool.Pool) *TxManager {
	return &TxManager{
		defaultDb: stdlib.OpenDBFromPool(db),
	}
}

func (txm *TxManager) Do(
	ctx context.Context,
	f func(db qrm.DB) error,
) error {
	tx, err := txm.defaultDb.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(
			errors.New("txm.defaultDb.Begin: "),
			ErrTx,
			err,
		)
	}

	defer tx.Rollback() //nolint

	err = f(tx)

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
