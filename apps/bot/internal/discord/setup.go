package discord

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) SetupHandlers() {
	b.session.AddHandler(b.onInteractionCreate)
}

func (b *Bot) SetupIntents() {
	b.session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
}

func (b *Bot) SetupScommands() {
	b.interactions.Polls.Register(b.session)
}
