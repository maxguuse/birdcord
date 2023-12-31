package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/poll"
	"github.com/maxguuse/birdcord/libs/pubsub"
	"go.uber.org/fx"
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Callback() func(i *discordgo.Interaction)
	Autocomplete() (func(i *discordgo.Interaction), bool)
}

type Handler struct {
	Commands []Command
	Session  *discordgo.Session
	Pubsub   pubsub.PubSub
}

type HandlerOpts struct {
	fx.In

	Pubsub  pubsub.PubSub
	Session *discordgo.Session

	PollHandler *poll.Handler
}

func New(opts HandlerOpts) *Handler {
	return &Handler{
		Commands: []Command{
			opts.PollHandler,
		},
		Session: opts.Session,
		Pubsub:  opts.Pubsub,
	}
}

func (h *Handler) Register() error {
	discordCommands := make([]*discordgo.ApplicationCommand, 0, len(h.Commands))

	for _, cmd := range h.Commands {
		discordCommands = append(discordCommands, cmd.Command())

		go h.Pubsub.Subscribe(cmd.Command().Name+":command", cmd.Callback())

		if callback, exists := cmd.Autocomplete(); exists {
			go h.Pubsub.Subscribe(cmd.Command().Name+":autocomplete", callback)
		}
	}

	_, err := h.Session.ApplicationCommandBulkOverwrite(h.Session.State.User.ID, "", discordCommands)
	if err != nil {
		return err
	}

	return nil
}

var NewFx = fx.Options(
	poll.NewFx,

	fx.Provide(New),
)
