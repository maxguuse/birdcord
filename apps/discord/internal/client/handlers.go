package client

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func (c *Client) registerHandlers() {
	c.AddHandler(c.onInteractionCreate)
	c.AddHandler(c.onConnect)
	c.AddHandler(c.onDisconnect)
	c.AddHandler(c.onReady)
}

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

	c.CommandsHandler.Subscribe()
	cmds := c.CommandsHandler.GetCommands()

	if _, err := c.ApplicationCommandBulkOverwrite(c.State.User.ID, "", cmds); err != nil {
		c.Log.Error(
			"Error creating commands",
			slog.String("error", err.Error()),
		)
	}
}

func (c *Client) onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	c.Log.Info("Bot is disconnected!")
}

func (c *Client) onReady(_ *discordgo.Session, _ *discordgo.Ready) {
	c.Log.Info("Bot is ready!")
}
