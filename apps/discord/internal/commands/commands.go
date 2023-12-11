package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/apps/discord/internal/postgres"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/samber/lo"
)

type command struct {
	Command      *discordgo.ApplicationCommand
	Callback     eventbus.EventHandler
	Autocomplete eventbus.EventHandler
}

type Handler struct {
	commands []*command
	eventbus *eventbus.EventBus

	Log      logger.Logger
	Database *postgres.Postgres
}

func New(
	eb *eventbus.EventBus,
	log logger.Logger,
	db *postgres.Postgres,
) *Handler {
	return &Handler{
		commands: commands,
		eventbus: eb,
		Log:      log,
		Database: db,
	}
}

func (h *Handler) GetCommands() []*discordgo.ApplicationCommand {
	commandsList := lo.Map(h.commands, func(c *command, _ int) *discordgo.ApplicationCommand {
		return c.Command
	})

	return commandsList
}

func (h *Handler) Subscribe() {
	for _, command := range commands {
		if command.Callback != nil {
			go h.eventbus.Subscribe(command.Command.Name+":command", command.Callback)
		}

		if command.Autocomplete != nil {
			go h.eventbus.Subscribe(command.Command.Name+":autocomplete", command.Autocomplete)
		}
	}
}

var commands = []*command{
	{
		Command:      poll,
		Callback:     &PollCommandHandler{},
		Autocomplete: nil,
	},
}

var poll = &discordgo.ApplicationCommand{
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
	},
}
