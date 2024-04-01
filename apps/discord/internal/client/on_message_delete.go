package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) onMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	ctx := context.Background()

	err := c.db.Messages().DeleteMessage(ctx, m.ID)
	if err != nil {
		c.logger.Error("error deleting message", slog.String("error", err.Error()))

		return
	}
	c.logger.Info(
		"message deleted",
		slog.String("message_id", m.ID),
	)
}
