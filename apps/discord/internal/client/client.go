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

type ClientOpts struct {
	fx.In
	LC fx.Lifecycle

	Log      logger.Logger
	Database *db.DB
	EB       *eventbus.EventBus
	CH       *commands.Handler
	Cfg      *config.Config
	S        *discordgo.Session
}

func New(opts ClientOpts) {
	client := &Client{
		Log:             opts.Log,
		Database:        opts.Database,
		Eventbus:        opts.EB,
		CommandsHandler: opts.CH,
		Session:         opts.S,
	}

	client.registerLogger()
	client.registerHandlers()

	opts.LC.Append(fx.Hook{
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
