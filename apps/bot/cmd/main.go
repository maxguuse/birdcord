package main

import (
	"github.com/maxguuse/birdcord/apps/bot/internal/discord"
	"github.com/maxguuse/birdcord/apps/bot/internal/interactions"
	"github.com/maxguuse/birdcord/libs/grpc/clients"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		interactions.NewFx,

		fx.Provide(
			clients.NewPolls,
		),

		fx.Invoke(discord.New),
	).Run()
}
