package poll

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func buildPollEmbed(
	poll *domain.PollWithDetails,
	user *discordgo.User,
) []*discordgo.MessageEmbed {
	optionsList := lo.Map(poll.Options, func(option domain.PollOption, i int) string {
		return fmt.Sprintf("**%d**. %s", i+1, option.Title)
	})

	return []*discordgo.MessageEmbed{
		{
			Title:       poll.Title,
			Description: strings.Join(optionsList, "\n"),
			Timestamp:   poll.CreatedAt.Format(time.RFC3339),
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
					Value:  strconv.Itoa(len(poll.Votes)),
					Inline: true,
				},
			},
		},
	}
}
