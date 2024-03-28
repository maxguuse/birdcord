package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
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
	CommandPoll            = "poll"
	SubcommandStart        = "start"
	SubcommandStop         = "stop"
	SubcommandStatus       = "status"
	SubcommandAddOption    = "add-option"
	SubcommandRemoveOption = "remove-option"
)

type optionsMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

type Handler struct {
	Log     logger.Logger
	Session *discordgo.Session

	// Refactor
	service *service.Service
}

type HandlerOpts struct {
	fx.In

	Log     logger.Logger
	Session *discordgo.Session

	// Refactor
	Service *service.Service
}

func NewHandler(opts HandlerOpts) *Handler {
	h := &Handler{
		Log:     opts.Log,
		Session: opts.Session,

		// Refactor
		service: opts.Service,
	}

	return h
}
