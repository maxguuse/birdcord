package embeds

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

var (
	ActivePoll = func(
		title, description, pollId string,
		totalVotes int32,
		nickname, avatarURL string,
	) *discordgo.MessageEmbed {
		return &discordgo.MessageEmbed{
			Title:       title,
			Description: description,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Color:       0x4d58d3,
			Type:        discordgo.EmbedTypeRich,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    nickname,
				IconURL: avatarURL,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: pollId,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Всего голосов",
					Value:  strconv.Itoa(int(totalVotes)),
					Inline: false,
				},
			},
		}
	}

	PollResults = func() *discordgo.MessageEmbed {

		return nil
	}
)
