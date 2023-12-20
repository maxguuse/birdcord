package main

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/client"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,

		commands.NewFx,

		fx.Provide(
			config.New,
			logger.New("discord"),
			db.New,
			eventbus.New,
		),
		fx.Invoke(
			client.New,
		),
	).Run()
}
