package poll

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func buildCommandOptionsMap(i *discordgo.Interaction) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	commandOptions := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		commandOptions[option.Name] = option
	}
	return commandOptions
}

func buildPollEmbed(
	poll queries.Poll,
	optionsList []string,
	user *discordgo.User,
	votesAmount int,
) []*discordgo.MessageEmbed {
	return []*discordgo.MessageEmbed{
		{
			Title:       poll.Title,
			Description: strings.Join(optionsList, "\n"),
			Timestamp:   poll.CreatedAt.Time.Format(time.RFC3339),
			Color:       0x4d58d3,
			Type:        discordgo.EmbedTypeRich,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.Username,
				IconURL: user.AvatarURL(""),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprint("Poll ID: ", poll.ID),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Всего голосов",
					Value:  strconv.Itoa(votesAmount),
					Inline: true,
				},
			},
		},
	}
}
