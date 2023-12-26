package client

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) onInteractionCreate(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		c.Log.Info("Got command",
			slog.String("command", i.ApplicationCommandData().Name),
			slog.String("user", i.Member.User.Username),
		)

		c.Pubsub.Publish(i.ApplicationCommandData().Name+":command", i.Interaction)
	case discordgo.InteractionApplicationCommandAutocomplete:
		c.Log.Info("Got autocomplete",
			slog.String("command", i.ApplicationCommandData().Name),
			slog.String("user", i.Member.User.Username),
		)

		c.Pubsub.Publish(i.ApplicationCommandData().Name+":autocomplete", i.Interaction)
	case discordgo.InteractionMessageComponent:
		c.Log.Info("Got message component",
			slog.String("custom_id", i.MessageComponentData().CustomID),
			slog.String("user", i.Member.User.Username),
		)

		c.Pubsub.Publish(i.MessageComponentData().CustomID, i.Interaction)
	}
}
