package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/scommands"
)

func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if i.ApplicationCommandData().Name == "poll" {
			scommands.PollCommandHandler(s, i)
			return
		}
	case discordgo.InteractionMessageComponent:

	}
}
