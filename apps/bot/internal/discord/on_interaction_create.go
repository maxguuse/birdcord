package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := b.interactions.Commands[i.ApplicationCommandData().Name]; ok {
			h(s, i.Interaction)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		fmt.Println("Got Autocompletion interaction, not implemented yet")
	case discordgo.InteractionMessageComponent:
		fmt.Println("Got Component interaction, not implemented yet")
	case discordgo.InteractionModalSubmit:
		fmt.Println("Got Modal Submit interaction, not implemented yet")
	case discordgo.InteractionPing:
		fmt.Println("Got Ping interaction, not implemented yet")
	}
}
