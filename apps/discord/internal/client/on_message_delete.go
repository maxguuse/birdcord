package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) onMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	ctx := context.Background()

	err := c.Database.Messages().DeleteMessage(ctx, m.ID)
	if err != nil {
		c.Log.Error("error deleting message", slog.String("error", err.Error()))

		return
	}
	c.Log.Info(
		"message deleted",
		slog.String("message_id", m.ID),
	)
}
