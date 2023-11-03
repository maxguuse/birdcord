package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) onReady(_ *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot is ready", r.User.Username, r.User.Discriminator)
}
