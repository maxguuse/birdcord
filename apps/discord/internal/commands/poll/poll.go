package poll

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

var NewFx = fx.Options(
	fx.Provide(
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
	Log      logger.Logger
	Database repository.DB
	Session  *discordgo.Session
}

type HandlerOpts struct {
	fx.In

	Log      logger.Logger
	Database repository.DB
	Session  *discordgo.Session
}

func NewHandler(opts HandlerOpts) *Handler {
	h := &Handler{
		Log:      opts.Log,
		Database: opts.Database,
		Session:  opts.Session,
	}

	return h
}

func (h *Handler) Autocomplete() (func(i *discordgo.Interaction), bool) {
	return func(i *discordgo.Interaction) {
		data := i.ApplicationCommandData()
		h.Log.Debug("data", slog.Any("data", data))

		commandOptions := helpers.BuildOptionsMap(i)

		switch i.ApplicationCommandData().Options[0].Name {
		case SubcommandStop:
			h.autocompletePollList(i, commandOptions)
		case SubcommandStatus:
			h.autocompletePollList(i, commandOptions)
		case SubcommandAddOption:
			h.autocompletePollList(i, commandOptions)
		case SubcommandRemoveOption:
			h.removeOptionAutocomplete(i, commandOptions)
		}
	}, true
}
