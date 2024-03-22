package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
)

func (h *Handler) statusPoll(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	pollId := options["poll"].IntValue()

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

	err = h.sendPollMessage(ctx, i, poll)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно отправлен.", nil
}
