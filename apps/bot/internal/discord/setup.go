package discord

import "github.com/bwmarrin/discordgo"

func (b *Bot) SetupHandlers() {
	b.session.AddHandler(b.onMessageCreate)
}

func (b *Bot) SetupIntents() {
	b.session.Identify.Intents = discordgo.IntentGuildMessages
}
