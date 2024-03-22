package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/samber/lo"
)

func (h *Handler) removePollOption(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	pollId := options["poll"].IntValue()
	optionId := options["option"].IntValue()

	var repoErr *repository.NotFoundError
	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if errors.As(err, &repoErr) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		return "", ErrNotAuthor
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return "", ErrNotFound
	}

	optionVotes := lo.CountBy(poll.Votes, func(v domain.PollVote) bool {
		return v.OptionID == int(optionId)
	})

	if optionVotes > 0 {
		return "", ErrOptionHasVotes
	}

	if len(poll.Options) <= 2 {
		return "", ErrTooFewOptions
	}

	err = h.Database.Polls().RemovePollOption(ctx, int(optionId))
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	poll.Options = lo.Filter(poll.Options, func(o domain.PollOption, _ int) bool {
		return o.ID != int(optionId)
	})

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Вариант опроса успешно удален.", nil
}
