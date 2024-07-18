package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/repository"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
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
	CommandPoll            = "poll"
	SubcommandStart        = "start"
	SubcommandStop         = "stop"
	SubcommandStatus       = "status"
	SubcommandAddOption    = "add-option"
	SubcommandRemoveOption = "remove-option"
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
		Name:         CommandPoll,
		Description:  "Управление опросами",
		DMPermission: lo.ToPtr(false),
	})

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.start)

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.stop).WithAutocompletion(h.autocompletePollList)

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.status).WithAutocompletion(h.autocompletePollList)

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.addOption).WithAutocompletion(h.autocompletePollList)

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.removeOption).WithAutocompletion(h.removeOptionAutocomplete)

	router.HandleComponent("poll-vote-btn", h.VoteBtnHandler)
}
