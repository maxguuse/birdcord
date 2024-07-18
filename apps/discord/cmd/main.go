package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/apps/discord/internal/client"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/redis/go-redis/v9"
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

			func(ctx context.Context, cfg *config.Config, lc fx.Lifecycle) *pgxpool.Pool {
				pool, err := pgxpool.New(ctx, cfg.ConnectionString)
				if err != nil {
					panic(err)
				}

				lc.Append(fx.StopHook(func() {
					pool.Close()
				}))

				return pool
			},

			func(cfg *config.Config) *redis.Client {
				return redis.NewClient(&redis.Options{
					Addr:     cfg.RedisDSN,
					Password: cfg.RedisPassword,
				})
			},

			txmanager.NewManager,
			txmanager.NewGetter,
		),

		modules.NewFx,

		client.NewFx,
	).Run()
}
