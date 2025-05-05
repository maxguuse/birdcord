package modules

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole"
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

	LiveroleHandler *liverole.Handler
}

func New(opts HandlerOpts) *Handler {
	return &Handler{
		Modules: []Module{
			opts.LiveroleHandler,
		},
	}
}

func (h *Handler) Register(router *disroute.Router) {
	for _, cmd := range h.Modules {
		cmd.Register(router)
	}
}

var NewFx = fx.Options(
	liverole.NewFx,

	fx.Provide(
		New,
	),
)
