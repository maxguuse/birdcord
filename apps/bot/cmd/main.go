package main

import (
	"github.com/maxguuse/birdcord/apps/bot/internal/discord"
	"github.com/maxguuse/birdcord/apps/bot/internal/scommands/polls"
	"github.com/maxguuse/birdcord/libs/grpc/clients"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,

		fx.Provide(
			clients.NewPolls,
			polls.New,
		),

		fx.Invoke(discord.New),
	).Run()
}
