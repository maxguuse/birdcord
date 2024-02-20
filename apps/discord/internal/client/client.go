package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

func NewSession(cfg *config.Config) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + cfg.DiscordToken)

	s.Identify.Intents = discordgo.IntentsAll

	return s, err
}

type Client struct {
	*discordgo.Session

	Cfg             *config.Config
	Log             logger.Logger
	Database        repository.DB
	CommandsHandler *commands.Handler
}

type ClientOpts struct {
	fx.In
	LC fx.Lifecycle

	Log      logger.Logger
	Database repository.DB
	Cfg      *config.Config

	Session         *discordgo.Session
	CommandsHandler *commands.Handler
}

func New(opts ClientOpts) *Client {
	client := &Client{
		Cfg:             opts.Cfg,
		Log:             opts.Log,
		Database:        opts.Database,
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
