package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/interactions/polls"
	"go.uber.org/fx"
)

type handlerMap map[string]func(s *discordgo.Session, i *discordgo.Interaction)

type Handlers struct {
	Polls *polls.Polls

	Commands     handlerMap
	Autocomplete handlerMap
	// MessageComponents
	// ModalSubmits
	// Pings
}

func NewHandlers(pollsHandlers *polls.Polls) *Handlers {
	return &Handlers{
		Polls: pollsHandlers,
		Commands: handlerMap{
			"poll": pollsHandlers.CommandHandler,
		},
		Autocomplete: handlerMap{
			"poll": pollsHandlers.AutocompleteHandler,
		},
	}
}

var NewFx = fx.Options(
	fx.Provide(
		polls.New,
		NewHandlers,
	),
)
