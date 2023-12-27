package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/pubsub"
	"go.uber.org/fx"
)

var NewFx = fx.Options(
	fx.Provide(
		NewVoteCallbackBuilder,
		NewHandler,
	),
)

type Handler struct {
	Log         logger.Logger
	Database    repository.DB
	Pubsub      pubsub.PubSub
	Session     *discordgo.Session
	VoteBuilder *VoteCallbackBuilder
}

type HandlerOpts struct {
	fx.In

	Log         logger.Logger
	Database    repository.DB
	Pubsub      pubsub.PubSub
	Session     *discordgo.Session
	VoteBuilder *VoteCallbackBuilder
}

func NewHandler(opts HandlerOpts) *Handler {
	return &Handler{
		Log:         opts.Log,
		Database:    opts.Database,
		Pubsub:      opts.Pubsub,
		Session:     opts.Session,
		VoteBuilder: opts.VoteBuilder,
	}
}

func (h *Handler) Command() *discordgo.ApplicationCommand {
	return command
}

func (h *Handler) Callback() func(i *discordgo.Interaction) {
	return func(i *discordgo.Interaction) {
		commandOptions := helpers.BuildOptionsMap(i)

		switch i.ApplicationCommandData().Options[0].Name {
		case "start":
			h.startPoll(i, commandOptions)
		case "stop":
			h.stopPoll(i, commandOptions)
		case "status":
			h.statusPoll(i, commandOptions)
		case "add-option":
			h.addPollOption(i, commandOptions)
		}
	}
}

func (h *Handler) Autocomplete() (func(i *discordgo.Interaction), bool) {
	return func(i *discordgo.Interaction) {
		commandOptions := helpers.BuildOptionsMap(i)

		switch i.ApplicationCommandData().Options[0].Name {
		case "stop":
			h.autocompletePollList(i, commandOptions)
		case "status":
			h.autocompletePollList(i, commandOptions)
		case "add-option":
			h.autocompletePollList(i, commandOptions)
		}
	}, true
}

var command = &discordgo.ApplicationCommand{
	Name:        "poll",
	Description: "Управление опросами",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "start",
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
			Name:        "stop",
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
			Name:        "status",
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
			Name:        "add-option",
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
	},
}
