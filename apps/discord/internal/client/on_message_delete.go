package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func (c *Client) onMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	ctx := context.Background()

	err := c.Database.Transaction(ctx, func(q *queries.Queries) error {
		msg, err := q.GetMessageByDiscordID(ctx, m.ID)
		if err != nil {
			return err
		}

		err = q.DeleteMessageById(ctx, msg.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.Log.Error("error deleting message", slog.String("error", err.Error()))

		return
	}
	c.Log.Info(
		"message deleted",
		slog.String("message_id", m.ID),
	)
}
