package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
)

type CommandHandler struct {
	Log      logger.Logger
	Database repository.DB
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

func NewCommandHandler(
	log logger.Logger,
	eb *eventbus.EventBus,
	db repository.DB,
	s *discordgo.Session,
) *CommandHandler {
	return &CommandHandler{
		Log:      log,
		Database: db,
		EventBus: eb,
		Session:  s,
	}
}

func (p *CommandHandler) Handle(i any) {
	cmd, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	commandOptions := buildCommandOptionsMap(cmd)

	switch cmd.ApplicationCommandData().Options[0].Name {
	case "start":
		p.startPoll(cmd, commandOptions)
	case "stop":
		p.stopPoll(cmd, commandOptions)
	}
}
