package liverole

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/pubsub"
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

type Handler struct {
	Log      logger.Logger
	Database repository.DB
	Pubsub   pubsub.PubSub
	Session  *discordgo.Session
}

type HandlerOpts struct {
	fx.In

	Log      logger.Logger
	Database repository.DB
	Pubsub   pubsub.PubSub
	Session  *discordgo.Session
}

func NewHandler(opts HandlerOpts) *Handler {
	return &Handler{
		Log:      opts.Log,
		Database: opts.Database,
		Pubsub:   opts.Pubsub,
		Session:  opts.Session,
	}
}

func (h *Handler) Command() *discordgo.ApplicationCommand {
	return command
}

func (h *Handler) Callback() func(i *discordgo.Interaction) {
	return func(i *discordgo.Interaction) {
		_ = helpers.BuildOptionsMap(i)
	}
}

func (h *Handler) Autocomplete() (func(i *discordgo.Interaction), bool) {
	return func(i *discordgo.Interaction) {
		data := i.ApplicationCommandData()
		h.Log.Debug("data", slog.Any("data", data))

		_ = helpers.BuildOptionsMap(i)
	}, false
}

var command = &discordgo.ApplicationCommand{
	Name:        "liverole",
	Description: "Управление live-ролями",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        SubcommandAdd,
			Description: "Добавить live-роль",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "role",
					Description:  "Роль",
					Type:         discordgo.ApplicationCommandOptionRole,
					Autocomplete: true,
					Required:     true,
				},
			},
		},
		{
			Name:        SubcommandRemove,
			Description: "Удалить live-роль",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "role",
					Description:  "Роль",
					Type:         discordgo.ApplicationCommandOptionRole,
					Autocomplete: true,
					Required:     true,
				},
			},
		},
		{
			Name:        SubcommandList,
			Description: "Список live-ролей",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        SubcommandClear,
			Description: "Очистить список live-ролей",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
	},
}
