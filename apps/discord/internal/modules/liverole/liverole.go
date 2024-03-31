package liverole

import (
	"github.com/bwmarrin/discordgo"
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
