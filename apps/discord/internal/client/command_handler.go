package client

import (
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"log/slog"
)

type Command interface {
	Execute(*Client, *discordgo.InteractionCreate)
	GetCommand() *discordgo.ApplicationCommand
}

var Commands = map[string]Command{
	"poll": &PollCommand{},
}

func (c *Client) RegisterCommands() error {
	commands := lo.MapToSlice(Commands, func(_ string, cmd Command) *discordgo.ApplicationCommand {
		return cmd.GetCommand()
	})

	_, err := c.ApplicationCommandBulkOverwrite(c.State.User.ID, "", commands)
	if err != nil {
		return err
	}

	return nil
}

type applicationCommandHandler struct{}

func (h *applicationCommandHandler) Handle(c *Client, i *discordgo.InteractionCreate) {
	c.Log.Debug(
		"Got application command",
		slog.String("command", i.ApplicationCommandData().Name),
	)

	cmd, exists := Commands[i.ApplicationCommandData().Name]

	if !exists {
		c.Log.Error(
			"Unknown command",
			slog.String("command", i.ApplicationCommandData().Name),
		)
		return
	}

	cmd.Execute(c, i)
}

type applicationCommandAutocompleteHandler struct{}

func (h *applicationCommandAutocompleteHandler) Handle(c *Client, i *discordgo.InteractionCreate) {
	c.Log.Debug(
		"Got application command autocomplete",
		slog.String("command", i.ApplicationCommandData().Name),
	)
}
