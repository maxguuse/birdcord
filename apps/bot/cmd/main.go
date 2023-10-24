package main

import (
	"github.com/maxguuse/birdcord/apps/bot/internal/discord"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Invoke(discord.New),
	).Run()
}
