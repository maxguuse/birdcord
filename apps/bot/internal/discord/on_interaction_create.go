package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println(i.Type)
	switch i.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		fallthrough
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "poll":
			b.polls.Handler(s, i)
		}
	case discordgo.InteractionMessageComponent:
		fmt.Println(i.MessageComponentData().CustomID) // poll_25_choice_72
		if strings.HasPrefix(i.MessageComponentData().CustomID, "poll") {
			b.polls.Handler(s, i)
		}
	}
}
