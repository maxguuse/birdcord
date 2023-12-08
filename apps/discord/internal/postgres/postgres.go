package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/libs/config"
)

type Postgres struct {
	Pool         *pgxpool.Pool
	QueryBuilder squirrel.StatementBuilderType
}

func New(cfg *config.Config) (*Postgres, error) {
	//TODO move to libs
	db, err := pgxpool.New(context.Background(), cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Postgres{
		Pool:         db,
		QueryBuilder: queryBuilder,
	}, nil
}

func (p *Postgres) Transaction(f func(pgx.Tx) error) error {
	tx, err := p.Pool.Begin(context.Background())
	if err != nil {
		return err
	}

	err = f(tx)

	if err != nil {
		_ = tx.Rollback(context.Background())

		return err
	}

	return tx.Commit(context.Background())
}
