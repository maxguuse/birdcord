package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

type CommandHandler struct {
	Log      logger.Logger
	Database repository.DB
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

type CommandHandlerOpts struct {
	fx.In

	Log      logger.Logger
	Database repository.DB
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

func NewCommandHandler(opts CommandHandlerOpts) *CommandHandler {
	return &CommandHandler{
		Log:      opts.Log,
		Database: opts.Database,
		EventBus: opts.EventBus,
		Session:  opts.Session,
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
