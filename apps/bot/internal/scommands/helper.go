package scommands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func buildOptionsMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionsMap[option.Name] = option
	}
	return optionsMap
}

func buildPollMessage(title string, options []string, user *discordgo.User) *discordgo.InteractionResponseData {
	// ! Pass all custom ids to the function
	var description string
	buttons := make([]discordgo.MessageComponent, 0)

	for i, option := range options {
		description += fmt.Sprintf("%d %s \n", i+1, option)
		buttons = append(buttons, discordgo.Button{
			Label:    option,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("%d", i), // ! TODO: CustomID format should be poll_<poll_id>_choice_<choice_id>
		})
	}

	actionRows := make([]discordgo.MessageComponent, (len(buttons)+4)/5)
	for i := 0; i < len(buttons); i += 5 {
		actionRow := discordgo.ActionsRow{}
		for j := 0; j < 5; j++ {
			if i+j >= len(buttons) {
				break
			}
			actionRow.Components = append(actionRow.Components, buttons[i+j])
		}
		actionRows[i/5] = actionRow
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,       // Title of the poll
		Description: description, // Options will be here
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Color:       0x4d58d3,
		Type:        discordgo.EmbedTypeRich,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL("1024"),
		},
	}

	responseData := &discordgo.InteractionResponseData{
		Content: "Опрос успешно создан!",
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
		Components: actionRows,
	}
	return responseData
}
