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
		return "", ErrNotAuthor
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return "", ErrNotFound
	}

	var maxVotes int = 0
	for _, option := range poll.Options {
		optionVotes := lo.CountBy(poll.Votes, func(v domain.PollVote) bool {
			return v.OptionID == option.ID
		})

		optionsWithVotes[option] = optionVotes

		if optionVotes > maxVotes {
			maxVotes = optionVotes
		}
	}

	winners := lo.FilterMap(poll.Options, func(o domain.PollOption, _ int) (string, bool) {
		return o.Title, optionsWithVotes[o] == maxVotes
	})

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
		stop:        true,
		fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Победители",
				Value:  strings.Join(winners, ","),
				Inline: true,
			},
		},
	})
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	err = h.Database.Polls().UpdatePollStatus(ctx, int(pollId), false)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно остановлен.", nil
}
