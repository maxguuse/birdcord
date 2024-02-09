package client

import (
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

	customStatus := lo.If(
		c.Cfg.Environment == "prod",
		"Released",
	).Else("Смотрит как Гусь кодит 💻")

	if err := c.UpdateStatusComplex(discordgo.UpdateStatusData{
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
		c.Log.Error(
			"Error updating status",
			slog.String("error", err.Error()),
		)
	}

	// guildsIds := lo.Map(r.Guilds, func(g *discordgo.Guild, _ int) string {
	// 	c.Log.Info(
	// 		"Connected guild",
	// 		slog.String("id", g.ID),
	// 		slog.String("name", g.Name),
	// 	)

	// 	return g.ID
	// })

	// newGuildsCount, err := c.Database.Queries().CreateGuilds(context.Background(), guildsIds)
	// if err != nil {
	// 	panic(err)
	// }

	// c.Log.Info(
	// 	"Created guilds",
	// 	slog.Int("new", int(newGuildsCount)),
	// )
}
