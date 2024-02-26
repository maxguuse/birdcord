package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/liverole"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/poll"
	"github.com/maxguuse/disroute"
	"go.uber.org/fx"
)

type Command interface {
	GetRoutes() *disroute.Cmd
	GetDiscordGo() *discordgo.ApplicationCommand
}

type Handler struct {
	Commands []Command
	Session  *discordgo.Session
	Router   *disroute.Router
}

type HandlerOpts struct {
	fx.In

	Session *discordgo.Session
	Router  *disroute.Router

	PollHandler     *poll.Handler
	LiveroleHandler *liverole.Handler
}

func New(opts HandlerOpts) *Handler {
	return &Handler{
		Commands: []Command{
			opts.PollHandler,
			opts.LiveroleHandler,
		},
		Session: opts.Session,
	}
}

func (h *Handler) Register() error {
	routes := make([]*disroute.Cmd, 0, len(h.Commands))
	discordgo := make([]*discordgo.ApplicationCommand, 0, len(h.Commands))
	for _, cmd := range h.Commands {
		routes = append(routes, cmd.GetRoutes())
		discordgo = append(discordgo, cmd.GetDiscordGo())
	}

	err := h.Router.RegisterAll(routes)
	if err != nil {
		return err
	}

	_, err = h.Session.ApplicationCommandBulkOverwrite(h.Session.State.User.ID, "", discordgo)
	if err != nil {
		return err
	}

	return nil
}

var NewFx = fx.Options(
	poll.NewFx,
	liverole.NewFx,

	fx.Provide(
		disroute.New,
		New,
	),
)
