package main

import (
	"github.com/maxguuse/birdcord/apps/polls/internal/db"
	"github.com/maxguuse/birdcord/apps/polls/internal/grpc"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			db.New,
		),

		fx.Invoke(
			grpc.StartPollsServer,
		),
	).Run()
}
