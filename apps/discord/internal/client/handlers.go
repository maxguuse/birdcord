package client

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

type InteractionHandler interface {
	Handle(*Client, *discordgo.InteractionCreate)
}

var interactions = map[discordgo.InteractionType]InteractionHandler{
	discordgo.InteractionApplicationCommand:             &applicationCommandHandler{},
	discordgo.InteractionApplicationCommandAutocomplete: &applicationCommandAutocompleteHandler{},
	discordgo.InteractionMessageComponent:               &messageComponentHandler{},
}

func (c *Client) onInteractionCreate(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	c.Log.Debug(
		"Got interaction",
		slog.String("type", i.Type.String()),
	)

	handler, exists := interactions[i.Type]

	if !exists {
		c.Log.Error(
			"Unhandled interaction type",
			slog.String("type", i.Type.String()),
		)
		return
	}

	handler.Handle(c, i)
}

func (c *Client) onConnect(_ *discordgo.Session, _ *discordgo.Connect) {
	c.Log.Info("Bot is connected!")

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

	if err := c.RegisterCommands(); err != nil {
		c.Log.Error(
			"Error registering commands",
			slog.String("error", err.Error()),
		)
	}
}
