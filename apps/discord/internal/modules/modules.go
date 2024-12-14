package modules

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice"
	"github.com/maxguuse/disroute"
	"go.uber.org/fx"
)

type Module interface {
	Register(*disroute.Router)
}

type Handler struct {
	Modules []Module
}

type HandlerOpts struct {
	fx.In

	PollHandler      *poll.Handler
	LiveroleHandler  *liverole.Handler
	TempvoiceHandler *tempvoice.Handler
}

func New(opts HandlerOpts) *Handler {
	return &Handler{
		Modules: []Module{
			opts.PollHandler,
			opts.LiveroleHandler,
			opts.TempvoiceHandler,
		},
	}
}

func (h *Handler) Register(router *disroute.Router) {
	for _, cmd := range h.Modules {
		cmd.Register(router)
	}
}

var NewFx = fx.Options(
	poll.NewFx,
	liverole.NewFx,
	tempvoice.NewFx,

	fx.Provide(
		New,
	),
)
