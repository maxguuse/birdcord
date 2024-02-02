package client

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) registerHandlers() {
	c.AddHandler(c.onInteractionCreate)
	c.AddHandler(c.onMessageDelete)
	c.AddHandler(c.onReady)
	c.AddHandler(c.onConnect)
	c.AddHandler(c.onDisconnect)
	c.AddHandler(c.onStatusChanged)
}

func (c *Client) onConnect(_ *discordgo.Session, _ *discordgo.Connect) {
	c.Log.Info("Bot is connected!")
}

func (c *Client) onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	c.Log.Info("Bot is disconnected!")
}

func (c *Client) onStatusChanged(_ *discordgo.Session, u *discordgo.PresenceUpdate) {
	if u.User.Bot {
		return
	}

	c.Log.Debug("Status changed",
		slog.String("user", u.Presence.User.ID),
		slog.String("status", string(u.Status)),
	)
}
