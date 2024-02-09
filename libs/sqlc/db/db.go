package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

type DB struct {
	queries *queries.Queries
	pool    *pgxpool.Pool
}

func MustInit(cfg *config.Config) *DB {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return &DB{
		pool:    pool,
		queries: queries.New(pool),
	}
}

func (p *DB) Queries() *queries.Queries {
	return p.queries
}

var (
	ErrTransactionFailed = errors.New("transaction failed")
	ErrTxBegin           = errors.New("could not begin transaction")
	ErrTxCommit          = errors.New("could not commit transaction")
	ErrTxRollback        = errors.New("could not rollback transaction")
)

func (p *DB) Transaction(ctx context.Context, f func(*queries.Queries) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return errors.Join(ErrTxBegin, err)
	}

	q := p.queries.WithTx(tx)

	err = f(q)

	if err != nil {
		tErr := tx.Rollback(ctx)
		if tErr != nil {
			return errors.Join(ErrTxRollback, tErr, err)
		}

		return errors.Join(ErrTransactionFailed, err)
	}

	cErr := tx.Commit(ctx)
	if cErr != nil {
		return errors.Join(ErrTxCommit, cErr)
	}

	return nil
}
