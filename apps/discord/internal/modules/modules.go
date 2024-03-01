package modules

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll"
	"github.com/maxguuse/disroute"
	"go.uber.org/fx"
)

type Module interface {
	GetRoutes() *disroute.Cmd
	GetDiscordGo() *discordgo.ApplicationCommand
	GetComponents() *disroute.Component
}

type Handler struct {
	Modules []Module
	Session *discordgo.Session
	Router  *disroute.Router
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
		Modules: []Module{
			opts.PollHandler,
			opts.LiveroleHandler,
		},
		Router:  opts.Router,
		Session: opts.Session,
	}
}

func (h *Handler) Register() error {
	routes := make([]*disroute.Cmd, 0, len(h.Modules))
	discordgo := make([]*discordgo.ApplicationCommand, 0, len(h.Modules))
	components := make([]*disroute.Component, 0, len(h.Modules))
	for _, cmd := range h.Modules {
		routes = append(routes, cmd.GetRoutes())
		discordgo = append(discordgo, cmd.GetDiscordGo())
		tCmp := cmd.GetComponents()
		if tCmp.Key == "" || tCmp.Handler == nil {
			continue
		}
		components = append(components, cmd.GetComponents())
	}

	err := h.Router.RegisterAll(routes)
	if err != nil {
		return err
	}

	err = h.Router.RegisterComponents(components)
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
		func() *disroute.Router {
			return disroute.New(
				disroute.WithComponentFunc(func(ic *discordgo.InteractionCreate) (key string) {
					parts := strings.Split(ic.MessageComponentData().CustomID, ":")

					return parts[0]
				}),
			)
		},
		New,
	),
)
