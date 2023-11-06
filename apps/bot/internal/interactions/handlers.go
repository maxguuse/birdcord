package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/interactions/polls"
	"go.uber.org/fx"
)

type Handlers struct {
	Polls *polls.Polls

	Commands map[string]func(s *discordgo.Session, i *discordgo.Interaction)
	// Autocompletes
	// MessageComponents
	// ModalSubmits
	// Pings
}

func NewHandlers(pollsHandlers *polls.Polls) *Handlers {
	return &Handlers{
		Polls: pollsHandlers,
		Commands: map[string]func(s *discordgo.Session, i *discordgo.Interaction){
			"poll": pollsHandlers.CommandHandler,
		},
	}
}

var NewFx = fx.Options(
	fx.Provide(
		polls.New,
		NewHandlers,
	),
)
