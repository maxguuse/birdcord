package client

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

func (c *Client) onReady(_ *discordgo.Session, r *discordgo.Ready) {
	c.Log.Info(
		"Bot is ready",
		slog.String("user", r.User.Username),
		slog.String("session_id", r.SessionID),
	)

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

	cmds := c.CommandsHandler.GetCommands()

	if _, err := c.ApplicationCommandBulkOverwrite(c.State.User.ID, "", cmds); err != nil {
		c.Log.Error(
			"Error overwriting commands",
			slog.String("error", err.Error()),
		)
	}

	addGuildsQuery := c.Database.QueryBuilder.
		Insert("guilds").
		Columns("discord_guild_id").
		Suffix("ON CONFLICT (discord_guild_id) DO NOTHING")

	for _, g := range r.Guilds {
		c.Log.Info(
			"Connected guild",
			slog.String("id", g.ID),
			slog.String("name", g.Name),
		)
		addGuildsQuery = addGuildsQuery.Values(g.ID)
	}

	addGuildsSql, args, _ := addGuildsQuery.ToSql()
	c.Log.Debug("add guilds query", slog.String("query", addGuildsSql))

	err := c.Database.Transaction(func(tx pgx.Tx) error {
		_, execErr := tx.Exec(context.Background(), addGuildsSql, args...)

		return execErr
	})
	if err != nil {
		c.Log.Error(
			"Error syncing guilds to database",
			slog.String("error", err.Error()),
		)
	} else {
		c.Log.Info("Synced guilds to database")
	}
}
