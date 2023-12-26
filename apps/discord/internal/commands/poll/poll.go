package poll

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/pubsub"
	"github.com/samber/lo"
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
		commandOptions := buildCommandOptionsMap(i)

		switch i.ApplicationCommandData().Options[0].Name {
		case "start":
			h.startPoll(i, commandOptions)
		case "stop":
			h.stopPoll(i, commandOptions)
		case "status":
			h.statusPoll(i, commandOptions)
		}
	}
}

func (h *Handler) Autocomplete() (func(i *discordgo.Interaction), bool) {
	return func(i *discordgo.Interaction) {
		commandOptions := buildCommandOptionsMap(i)

		switch i.ApplicationCommandData().Options[0].Name {
		case "stop":
			h.autocompletePollList(i, commandOptions)
		case "status":
			h.autocompletePollList(i, commandOptions)
		}
	}, true
}

func (h *Handler) autocompletePollList(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	ctx := context.Background()

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return
	}

	polls, err := h.Database.Polls().GetActivePolls(ctx, guild.ID, user.ID)
	if err != nil {
		return
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(polls))
	for i, poll := range polls {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("Poll ID: %d | %s", poll.ID, poll.Title),
			Value: poll.ID,
		}
	}

	err = h.Session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: lo.Filter(choices, func(c *discordgo.ApplicationCommandOptionChoice, _ int) bool {
				s, ok := options["poll"].Value.(string)
				if !ok {
					return false
				}

				return strings.Contains(c.Name, s)
			}),
		},
	})
	if err != nil {
		return
	}
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
	},
}
