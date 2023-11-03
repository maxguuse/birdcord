package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"go.uber.org/fx"
	"os"
)

func New(lc fx.Lifecycle) queries.DBTX {
	conn, err := pgxpool.New(context.Background(), os.Getenv("CONNECTION_STRING"))
	if err != nil {
		fmt.Println("Error establishing connection with database:", err)
		panic(err)
	}

	lc.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			conn.Close()
			return nil
		},
	})

	return conn
}
