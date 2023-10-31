package polls

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

func buildOptionsMap(
	options []*discordgo.ApplicationCommandInteractionDataOption,
) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionsMap[option.Name] = option
	}
	return optionsMap
}

func buildPollEmbed(
	title string,
	description string,
	user *discordgo.Member,
	footer *discordgo.MessageEmbedFooter,
	fields []*discordgo.MessageEmbedField,
) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Color:       0x4d58d3,
		Type:        discordgo.EmbedTypeRich,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Nick,
			IconURL: user.AvatarURL("1024"),
		},
		Footer: footer,
		Fields: fields,
	}

	return embed
}

func parseVoteFromButtonInteraction(i *discordgo.InteractionCreate) *Vote {
	customId := i.MessageComponentData().CustomID
	customIdParts := strings.Split(customId, "_")

	if len(customIdParts) != 4 {
		fmt.Println("Error parsing CustomID: ", customId, "len(customIdParts) != 4, invalid format")
		return nil
	}

	pollId, err := strconv.Atoi(customIdParts[1])
	optionId, err := strconv.Atoi(customIdParts[3])
	if err != nil {
		fmt.Println("Error parsing CustomID: ", customId, err)
		return nil
	}

	return &Vote{
		PollID:   int32(pollId),
		OptionID: int32(optionId),
		UserID:   i.Member.User.ID,
	}
}
