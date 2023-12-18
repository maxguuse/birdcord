package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

type DB struct {
	queries *queries.Queries
	pool    *pgxpool.Pool
}

func New(cfg *config.Config) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	return &DB{
		pool:    pool,
		queries: queries.New(pool),
	}, nil
}

func (p *DB) Queries() *queries.Queries {
	return p.queries
}

func (p *DB) Transaction(f func(*queries.Queries) error) (transactionErr error, callbackError error) {
	ctx := context.Background()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err, nil
	}

	q := p.queries.WithTx(tx)

	err = f(q)

	if err != nil {
		_ = tx.Rollback(ctx)
	}

	return tx.Commit(ctx), err
}
