package liverole

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

var NewFx = fx.Options(
	fx.Provide(
		service.New,

		NewHandler,
	),
)

const (
	SubcommandAdd    = "add"
	SubcommandRemove = "remove"
	SubcommandList   = "list"
	SubcommandClear  = "clear"
)

type optionsMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

type Handler struct {
	Log     logger.Logger
	Session *discordgo.Session

	Service *service.Service
}

type HandlerOpts struct {
	fx.In

	Log     logger.Logger
	Session *discordgo.Session

	service *service.Service
}

func NewHandler(opts HandlerOpts) *Handler {
	h := &Handler{
		Log:     opts.Log,
		Session: opts.Session,
		Service: opts.service,
	}

	return h
}
