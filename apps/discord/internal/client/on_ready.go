package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (c *Client) onReady(_ *discordgo.Session, r *discordgo.Ready) {
	c.Log.Info(
		"Bot is ready",
		slog.String("user", r.User.Username),
		slog.String("session_id", r.SessionID),
	)

	err := c.CommandsHandler.Register()
	if err != nil {
		panic(err)
	}

	if err := c.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: string(discordgo.StatusOnline),
		AFK:    false,
		Activities: []*discordgo.Activity{
			{
				Name: "как Гусь кодит 💻",
				Type: discordgo.ActivityTypeWatching,
			},
		},
	}); err != nil {
		c.Log.Error(
			"Error updating status",
			slog.String("error", err.Error()),
		)
	}

	guildsIds := lo.Map(r.Guilds, func(g *discordgo.Guild, _ int) string {
		c.Log.Info(
			"Connected guild",
			slog.String("id", g.ID),
			slog.String("name", g.Name),
		)

		return g.ID
	})
	newGuildsCount, err := c.Database.Queries().CreateGuilds(context.Background(), guildsIds)
	if err != nil {
		c.Log.Error(
			"Error creating guilds",
			slog.String("error", err.Error()),
		)
	} else {
		c.Log.Info(
			"Created guilds",
			slog.Int("new", int(newGuildsCount)),
		)
	}
}
