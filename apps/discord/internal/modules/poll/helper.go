package poll

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

const (
	VOTES_BAR_BLOCK = "■"
	VOTES_BAR_SPACE = " "
)

func buildPollEmbed(
	poll *domain.PollWithDetails,
	user *discordgo.User,
) []*discordgo.MessageEmbed {
	totalVotes := len(poll.Votes)

	optionsList := lo.Map(poll.Options, func(option domain.PollOption, i int) string {
		return fmt.Sprintf("**%d**. %s", i+1, option.Title)
	})

	optionsPercentageBars := lo.Map(poll.Options, func(option domain.PollOption, i int) string {
		votesForOption := lo.CountBy(poll.Votes, func(vote domain.PollVote) bool {
			return vote.OptionID == option.ID
		})

		percentage := (float64(votesForOption) / float64(totalVotes)) * 100
		if math.IsNaN(percentage) {
			percentage = 0
		}

		t := math.Ceil(percentage)
		t2 := int(math.Floor(t / 3.33))

		if t2 < 0 {
			t2 = 0
		}

		bar := strings.Repeat(VOTES_BAR_BLOCK, t2) + strings.Repeat(VOTES_BAR_SPACE, 30-t2)

		return fmt.Sprintf("(%d) | %s | (%d%%)", i+1, bar, int(t))
	})

	optionsListDesc := strings.Join(optionsList, "\n")
	optionsBarsDesc := strings.Join(optionsPercentageBars, "\n")

	return []*discordgo.MessageEmbed{
		{
			Title:       poll.Title,
			Description: optionsListDesc + "\n```" + optionsBarsDesc + "```",
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
					Value:  strconv.Itoa(totalVotes),
					Inline: true,
				},
			},
		},
	}
}
