package client

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
	"log/slog"
)

type Client struct {
	*discordgo.Session

	Log             logger.Logger
	Eventbus        *eventbus.EventBus
	CommandsHandler *commands.Handler
}

func New(
	lc fx.Lifecycle,
	log logger.Logger,
	eb *eventbus.EventBus,
	ch *commands.Handler,
	cfg *config.Config,
) {
	client := &Client{
		Log:             log,
		Eventbus:        eb,
		CommandsHandler: ch,
	}

	if s, err := discordgo.New("Bot " + cfg.DiscordToken); err != nil {
		log.Error(
			"Error creating Discord session",
			slog.String("error", err.Error()),
		)
		return
	} else {
		client.Session = s
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
