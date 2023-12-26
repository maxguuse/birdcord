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

func New(cfg *config.Config) (*DB, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return &DB{
		pool:    pool,
		queries: queries.New(pool),
	}, nil
}

func (p *DB) Queries() *queries.Queries {
	return p.queries
}

func (p *DB) Transaction(f func(*queries.Queries) error) error {
	ctx := context.Background()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := p.queries.WithTx(tx)

	err = f(q)

	if err != nil {
		return errors.Join(tx.Rollback(ctx), err)
	}

	return tx.Commit(ctx)
}
