package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
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
	Database *db.DB
}

func New(
	eb *eventbus.EventBus,
	log logger.Logger,
	db *db.DB,
) *Handler {
	h := &Handler{
		commands: []*command{
			{
				Command: poll,
				Callback: &PollCommandHandler{
					Log:      log,
					Database: db,
				},
				Autocomplete: nil,
			},
		},
		eventbus: eb,
		Log:      log,
		Database: db,
	}

	for _, cmd := range h.commands {
		if cmd.Callback != nil {
			go h.eventbus.Subscribe(cmd.Command.Name+":command", cmd.Callback)
		}

		if cmd.Autocomplete != nil {
			go h.eventbus.Subscribe(cmd.Command.Name+":autocomplete", cmd.Autocomplete)
		}
	}

	return h
}

func (h *Handler) GetCommands() []*discordgo.ApplicationCommand {
	commandsList := lo.Map(h.commands, func(c *command, _ int) *discordgo.ApplicationCommand {
		return c.Command
	})

	return commandsList
}
