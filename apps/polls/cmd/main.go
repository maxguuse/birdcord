package main

import (
	"github.com/maxguuse/birdcord/apps/polls/internal/db"
	"github.com/maxguuse/birdcord/apps/polls/internal/grpc"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"go.uber.org/fx"
)

func main() {
	fx.New(

		fx.Provide(
			db.New,
			queries.New,
		),

		fx.Invoke(
			grpc.StartPollsServer,
		),
	).Run()
}
