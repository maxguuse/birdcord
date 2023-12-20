package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type command struct {
	Command      *discordgo.ApplicationCommand
	Callback     eventbus.EventHandler
	Autocomplete eventbus.EventHandler
}

type Handler struct {
	commands []*command
	eventbus *eventbus.EventBus

	pollCommandHandler *PollCommandHandler
}

func New(
	eb *eventbus.EventBus,
	pollCommandHandler *PollCommandHandler,
) *Handler {
	h := &Handler{
		commands: []*command{
			{
				Command:      poll,
				Callback:     pollCommandHandler,
				Autocomplete: nil,
			},
		},
		eventbus: eb,
	}

	for _, cmd := range h.commands {
		if cmd.Callback != nil {
			go h.eventbus.Subscribe(cmd.Command.Name+":command", cmd.Callback)
		}

		if cmd.Autocomplete != nil {
			go h.eventbus.Subscribe(cmd.Command.Name+":autocomplete", cmd.Autocomplete)
		}
	}

	return h
}

func (h *Handler) GetCommands() []*discordgo.ApplicationCommand {
	commandsList := lo.Map(h.commands, func(c *command, _ int) *discordgo.ApplicationCommand {
		return c.Command
	})

	return commandsList
}

var NewFx = fx.Options(
	fx.Provide(
		newPolls,
		New,
	),
)
