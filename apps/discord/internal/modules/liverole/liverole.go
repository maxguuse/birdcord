package liverole

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/disroute"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

var NewFx = fx.Options(
	fx.Provide(
		service.New,

		NewHandler,
	),
)

const (
	CommandLiverole  = "liverole"
	SubcommandAdd    = "add"
	SubcommandRemove = "remove"
	SubcommandList   = "list"
	SubcommandClear  = "clear"
)

type Handler struct {
	logger  logger.Logger
	Service *service.Service
}

type HandlerOpts struct {
	fx.In

	Log     logger.Logger
	Service *service.Service
}

func NewHandler(opts HandlerOpts) *Handler {
	h := &Handler{
		logger:  opts.Log,
		Service: opts.Service,
	}

	return h
}

func (h *Handler) Register(router *disroute.Router) {
	r := router.Mount(&discordgo.ApplicationCommand{
		Name:         CommandLiverole,
		Description:  "Управление live-ролями",
		DMPermission: lo.ToPtr(false),
	})

	r.Use(func(hf disroute.HandlerFunc) disroute.HandlerFunc {
		return func(ctx *disroute.Ctx) disroute.Response {
			resp := hf(ctx)

			if resp.Err != nil || resp.CustomResponse != nil {
				return resp
			}

			return disroute.Response{
				CustomResponse: &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: resp.Message,
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				},
			}
		}
	})

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.addLiveRole)

	r.Handle(&discordgo.ApplicationCommandOption{
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
	}, h.removeLiveRole)

	r.Handle(&discordgo.ApplicationCommandOption{
		Name:        SubcommandList,
		Description: "Список live-ролей",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	}, h.listLiveRoles)

	r.Handle(&discordgo.ApplicationCommandOption{
		Name:        SubcommandClear,
		Description: "Очистить список live-ролей",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	}, h.clearLiveRoles)
}
