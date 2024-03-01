package liverole

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/helpers"
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
	SubcommandAdd    = "add"
	SubcommandRemove = "remove"
	SubcommandList   = "list"
	SubcommandClear  = "clear"
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

		_ = helpers.BuildOptionsMap(i)
	}, false
}
