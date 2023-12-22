package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"go.uber.org/fx"
)

type Client struct {
	*discordgo.Session

	Log             logger.Logger
	Database        *db.DB
	Eventbus        *eventbus.EventBus
	CommandsHandler *commands.Handler
}

func New(
	lc fx.Lifecycle,
	log logger.Logger,
	db *db.DB,
	eb *eventbus.EventBus,
	ch *commands.Handler,
	cfg *config.Config,
	s *discordgo.Session,
) {
	client := &Client{
		Log:             log,
		Database:        db,
		Eventbus:        eb,
		CommandsHandler: ch,
		Session:         s,
	}

	client.registerLogger()
	client.registerHandlers()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := client.Open(); err != nil {
				client.Log.Error(
					"Error opening connection",
					slog.String("error", err.Error()),
				)
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := client.Close(); err != nil {
				client.Log.Error(
					"Error closing connection",
					slog.String("error", err.Error()),
				)
				return err
			}
			return nil
		},
	})
}
