package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/logger"
	"log/slog"
)

type PollCommandHandler struct {
	Log logger.Logger
}

func (p *PollCommandHandler) Handle(s *discordgo.Session, i interface{}) {
	commandData, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	p.Log.Info("poll command", slog.String("command", commandData.ApplicationCommandData().Name))
}
