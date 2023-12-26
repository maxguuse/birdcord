package main

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/client"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/pubsub"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			config.New,
			logger.New("discord"),
			fx.Annotate(
				pubsub.New(100),
				fx.As(new(pubsub.PubSub)),
			),
		),

		repository.NewFx,
		commands.NewFx,

		client.NewFx,
	).Run()
}
