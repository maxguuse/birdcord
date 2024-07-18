package client

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Client) registerHandlers() {
	c.router.Session().AddHandler(c.router.InteractionHandler)

	c.router.Session().AddHandler(c.onReady)
	c.router.Session().AddHandler(c.onConnect)
	c.router.Session().AddHandler(c.onDisconnect)
	c.router.Session().AddHandler(c.onStatusChanged)
}

func (c *Client) onConnect(_ *discordgo.Session, _ *discordgo.Connect) {
	c.logger.Info("Bot is connected!")
}

func (c *Client) onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	c.logger.Info("Bot is disconnected!")
}
