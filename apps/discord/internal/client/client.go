package client

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/postgres"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
	"log/slog"
)

type Client struct {
	*discordgo.Session

	Log      logger.Logger
	Database *postgres.Postgres
}

func New(lc fx.Lifecycle, log logger.Logger, postgres *postgres.Postgres, cfg *config.Config) {
	client := &Client{
		Log:      log,
		Database: postgres,
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

	client.LogLevel = discordgo.LogWarning
	discordgo.Logger = client.logger()

	client.AddHandler(client.onConnect)
	client.AddHandler(func(_ *discordgo.Session, _ *discordgo.Disconnect) {
		client.Log.Info("Bot is disconnected!")
	})
	client.AddHandler(func(_ *discordgo.Session, r *discordgo.Ready) {
		client.Log.Info("Bot is ready!")
	})

	client.AddHandler(client.onInteractionCreate)

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

type messageComponentHandler struct{}

func (h *messageComponentHandler) Handle(c *Client, i *discordgo.InteractionCreate) {
	c.Log.Debug(
		"Got message component",
		slog.Uint64("type", uint64(i.MessageComponentData().ComponentType)),
		slog.String("custom_id", i.MessageComponentData().CustomID),
	)
}
