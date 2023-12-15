package client

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func (c *Client) onInteractionCreate(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	c.Log.Debug(
		"Got interaction",
		slog.String("type", i.Type.String()),
	)

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		c.Eventbus.Publish(i.ApplicationCommandData().Name+":command", c.Session, i.Interaction)
	case discordgo.InteractionApplicationCommandAutocomplete:
		c.Eventbus.Publish(i.ApplicationCommandData().Name+":autocomplete", c.Session, i.Interaction)
	}
}
