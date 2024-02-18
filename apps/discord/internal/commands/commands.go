package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/poll"
	"github.com/maxguuse/birdcord/libs/pubsub"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
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
	discordCommands := lo.Map(h.Commands, func(cmd Command, _ int) *discordgo.ApplicationCommand {
		return cmd.Command()
	})

	wg := new(errgroup.Group)
	for _, cmd := range h.Commands {
		cmd := cmd

		wg.Go(func() error {
			return h.Pubsub.Subscribe(cmd.Command().Name+":command", cmd.Callback())
		})

		if callback, exists := cmd.Autocomplete(); exists {
			wg.Go(func() error {
				return h.Pubsub.Subscribe(cmd.Command().Name+":autocomplete", callback)
			})
		}
	}
	if err := wg.Wait(); err != nil {
		return err
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
