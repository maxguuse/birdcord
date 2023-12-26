package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/pubsub"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"go.uber.org/fx"
)

func NewSession(cfg *config.Config) (*discordgo.Session, error) {
	return discordgo.New("Bot " + cfg.DiscordToken)
}

type Client struct {
	*discordgo.Session

	Log             logger.Logger
	Database        *db.DB
	Pubsub          pubsub.PubSub
	CommandsHandler *commands.Handler
}

type ClientOpts struct {
	fx.In
	LC fx.Lifecycle

	Log             logger.Logger
	Database        *db.DB
	Pubsub          pubsub.PubSub
	Cfg             *config.Config
	Session         *discordgo.Session
	CommandsHandler *commands.Handler
}

func New(opts ClientOpts) *Client {
	client := &Client{
		Log:             opts.Log,
		Database:        opts.Database,
		Pubsub:          opts.Pubsub,
		Session:         opts.Session,
		CommandsHandler: opts.CommandsHandler,
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

	return client
}

var NewFx = fx.Options(
	fx.Provide(
		NewSession,
	),
	fx.Invoke(
		New,
	),
)
