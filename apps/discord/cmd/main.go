package main

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/client"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/apps/discord/internal/session"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,

		fx.Provide(
			config.New,
			logger.New("discord"),
			session.New,
			eventbus.New,
		),

		repository.NewFx,
		commands.NewFx,

		fx.Invoke(
			client.New,
		),
	).Run()
}
