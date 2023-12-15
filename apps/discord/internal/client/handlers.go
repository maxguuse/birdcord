package client

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Client) registerHandlers() {
	c.AddHandler(c.onInteractionCreate)
	c.AddHandler(c.onReady)
	c.AddHandler(c.onConnect)
	c.AddHandler(c.onDisconnect)
}

func (c *Client) onConnect(_ *discordgo.Session, _ *discordgo.Connect) {
	c.Log.Info("Bot is connected!")
}

func (c *Client) onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	c.Log.Info("Bot is disconnected!")
}
