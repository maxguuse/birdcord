package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/scommands"
)

func (b *Bot) SetupHandlers() {
	b.session.AddHandler(b.onInteractionCreate)
	b.session.AddHandler(b.onReady)
}

func (b *Bot) SetupIntents() {
	b.session.Identify.Intents = discordgo.IntentGuildMessages
}

func (b *Bot) SetupScommands() {
	scommands.Register(b.session)
}
