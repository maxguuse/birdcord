package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, err := b.session.ChannelMessageSend(m.ChannelID, "pong")
		if err != nil {
			fmt.Println("Error sending pong:", err)
		}
	}

	if m.Content == "pong" {
		_, err := b.session.ChannelMessageSend(m.ChannelID, "ping")
		if err != nil {
			fmt.Println("Error sending ping:", err)
		}
	}

	if m.Content == "reply" {
		_, err := b.session.ChannelMessageSendReply(m.ChannelID, m.Author.Mention()+" reply", m.Message.Reference())
		if err != nil {
			fmt.Println("Error sending pong:", err)
		}
	}
}
