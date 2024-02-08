package poll

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
	Pubsub   pubsub.PubSub
	Session  *discordgo.Session

	subcommandsHandlers map[string]func(*discordgo.Interaction, optionsMap) error
}

type HandlerOpts struct {
	fx.In

	Log      logger.Logger
	Database repository.DB
	Pubsub   pubsub.PubSub
	Session  *discordgo.Session
}

func NewHandler(opts HandlerOpts) *Handler {
	h := &Handler{
		Log:      opts.Log,
		Database: opts.Database,
		Pubsub:   opts.Pubsub,
		Session:  opts.Session,
	}

	h.subcommandsHandlers = map[string]func(*discordgo.Interaction, optionsMap) error{
		SubcommandStart:        h.startPoll,
		SubcommandStop:         h.stopPoll,
		SubcommandStatus:       h.statusPoll,
		SubcommandAddOption:    h.addPollOption,
		SubcommandRemoveOption: h.removePollOption,
	}

	return h
}

func (h *Handler) Command() *discordgo.ApplicationCommand {
	return command
}

func (h *Handler) Callback() func(i *discordgo.Interaction) {
	return func(i *discordgo.Interaction) {
		commandOptions := helpers.BuildOptionsMap(i)

		sh := h.subcommandsHandlers[i.ApplicationCommandData().Options[0].Name]
		if sh == nil {
			return
		}

		err := sh(i, commandOptions)
		if err != nil {
			_ = helpers.InteractionResponseProcess(h.Session, i, "", err)
		}
	}
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

var command = &discordgo.ApplicationCommand{
	Name:        "poll",
	Description: "Управление опросами",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        SubcommandStart,
			Description: "Начать опрос",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "title",
					Description: "Заголовок опроса",
					Type:        discordgo.ApplicationCommandOptionString,
					MaxLength:   50,
					Required:    true,
				},
				{
					Name:        "options",
					Description: "Варианты ответа (разделите их символом '|')",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        SubcommandStop,
			Description: "Остановить опрос",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "poll",
					Description:  "Опрос",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        SubcommandStatus,
			Description: "Статус опроса",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "poll",
					Description:  "Опрос",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        SubcommandAddOption,
			Description: "Добавить вариант ответа к опросу",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "poll",
					Description:  "Опрос",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:        "option",
					Description: "Новый вариант ответа",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					MaxLength:   50,
				},
			},
		},
		{
			Name:        SubcommandRemoveOption,
			Description: "Удалить вариант ответа из опроса",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "poll",
					Description:  "Опрос",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "option",
					Description:  "Вариант ответа",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	},
}
