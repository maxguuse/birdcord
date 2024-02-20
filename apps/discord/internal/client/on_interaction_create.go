package client

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) onInteractionCreate(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		c.onCommand(i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		c.onAutocomplete(i)
	case discordgo.InteractionMessageComponent:
		c.onMessageComponent(i)
	}
}

func (c *Client) onCommand(i *discordgo.InteractionCreate) {
	c.Log.Info("Got command",
		slog.String("command", i.ApplicationCommandData().Name),
		slog.String("user", i.Member.User.Username),
	)

	err := c.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		c.Log.Error("error responding to interaction", slog.String("error", err.Error()))

		return
	}

	c.CommandsHandler.Router.FindAndExecute(i) //TODO handle error
}

func (c *Client) onAutocomplete(i *discordgo.InteractionCreate) {
	c.Log.Info("Got autocomplete",
		slog.String("command", i.ApplicationCommandData().Name),
		slog.String("user", i.Member.User.Username),
	)

	c.Pubsub.Publish(i.ApplicationCommandData().Name+":autocomplete", i.Interaction)
}

func (c *Client) onMessageComponent(i *discordgo.InteractionCreate) {
	c.Log.Info("Got message component",
		slog.String("custom_id", i.MessageComponentData().CustomID),
		slog.String("user", i.Member.User.Username),
	)

	err := c.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		c.Log.Error("error responding to interaction", slog.String("error", err.Error()))

		return
	}

	c.Pubsub.Publish(i.MessageComponentData().CustomID, i.Interaction)
}
