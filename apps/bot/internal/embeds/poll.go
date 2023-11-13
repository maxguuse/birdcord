package embeds

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

var (
	ActivePoll = func(
		title, description, pollId string,
		nickname, avatarURL string,
		totalVotes int32,
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
					Inline: true,
				},
			},
		}
	}

	PollResults = func(
		title, description, pollId string,
		nickname, avatarURL string,
		winners string,
		totalVotes int32,
	) *discordgo.MessageEmbed {
		embed := ActivePoll(title, description, pollId, nickname, avatarURL, totalVotes)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Победители",
			Value:  winners,
			Inline: true,
		})

		return embed
	}
)
