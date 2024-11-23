package tempvoice

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/repository"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/service"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/disroute"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

var NewFx = fx.Options(
	fx.Provide(
		fx.Annotate(repository.NewPgx, fx.As(new(repository.Repository))),
		service.New,

		NewHandler,
	),
)

const (
	Command         = "tempvoice"
	SubcommandSetup = "setup"
)

type Handler struct {
	logger  logger.Logger
	service *service.Service
}

type HandlerOpts struct {
	fx.In

	Log     logger.Logger
	Service *service.Service
}

func NewHandler(opts HandlerOpts) *Handler {
	h := &Handler{
		logger:  opts.Log,
		service: opts.Service,
	}

	return h
}

func (h *Handler) Register(router *disroute.Router) {
	r := router.Mount(&discordgo.ApplicationCommand{
		Name:         Command,
		Description:  "Управление временными каналами",
		DMPermission: lo.ToPtr(false),
	})

	r.Handle(&discordgo.ApplicationCommandOption{
		Name:        SubcommandSetup,
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Description: "Настройка хабов",
	}, h.setup)

	router.HandleComponent("create-tempvoice-hub-btn", func(c *disroute.Ctx) disroute.Response {
		h.logger.Error("not implemented")

		return disroute.Response{
			Err: errors.New("not implemented"),
		}
	})

	router.HandleComponent(
		"configure-tempvoice-hub-select-menu",
		func(c *disroute.Ctx) disroute.Response {
			h.logger.Error("not implemented")

			return disroute.Response{
				Err: errors.New("not implemented"),
			}
		},
	)
}
