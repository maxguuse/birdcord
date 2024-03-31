package client

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (c *Client) onReady(_ *discordgo.Session, r *discordgo.Ready) {
	c.logger.Info(
		"Bot is ready",
		slog.String("user", r.User.Username),
		slog.String("session_id", r.SessionID),
	)

	customStatus := lo.If(
		c.cfg.Environment == "prod",
		"Release "+c.cfg.Version,
	).Else("Смотрит как Гусь кодит 💻")

	if err := c.router.Session().UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name:  "Custom status",
				Type:  discordgo.ActivityTypeCustom,
				State: customStatus,
			},
		},
		AFK:    false,
		Status: string(discordgo.StatusOnline),
	}); err != nil {
		c.logger.Error(
			"Error updating status",
			slog.String("error", err.Error()),
		)
	}
}
