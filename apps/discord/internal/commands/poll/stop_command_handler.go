package poll

import (
	"context"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) stopPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) (string, error) {
	ctx := context.Background()

	optionsWithVotes := make(map[domain.PollOption]int)

	pollId := options["poll"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		return "", &domain.UsersideError{
			Msg: "Для изменения опроса нужно быть его автором.",
		}
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return "", &domain.UsersideError{
			Msg: "Опроса не существует.",
		}
	}

	var maxVotes int = 0
	for _, option := range poll.Options {
		optionVotes := lo.Filter(poll.Votes, func(v domain.PollVote, _ int) bool {
			return v.OptionID == option.ID
		})
		optionVotesAmount := len(optionVotes)

		optionsWithVotes[option] = optionVotesAmount

		if optionVotesAmount > maxVotes {
			maxVotes = optionVotesAmount
		}
	}

	winners := make([]domain.PollOption, 0, len(poll.Options))

	for _, option := range poll.Options {
		if optionsWithVotes[option] == maxVotes {
			winners = append(winners, option)
		}
	}

	winnersList := lo.Map(winners, func(option domain.PollOption, _ int) string {
		return option.Title
	})

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
		stop:        true,
		fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Победители",
				Value:  strings.Join(winnersList, ","),
				Inline: true,
			},
		},
	})

	err = h.Database.Polls().UpdatePollStatus(ctx, int(pollId), false)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно остановлен.", nil
}
