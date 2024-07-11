package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/apps/discord/internal/client"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			func() context.Context {
				return context.Background()
			},
			config.New,
			logger.New("discord"),

			func(ctx context.Context, cfx *config.Config, lc fx.Lifecycle) *pgxpool.Pool {
				pool, err := pgxpool.New(ctx, cfx.ConnectionString)
				if err != nil {
					panic(err)
				}

				lc.Append(fx.StopHook(func() {
					pool.Close()
				}))

				return pool
			},
			txmanager.New,
		),

		repository.NewFx,
		modules.NewFx,

		client.NewFx,
	).Run()
}
