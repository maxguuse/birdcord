package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) statusPoll(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.GetPoll(ctx, &service.GetPollRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  options["poll"].IntValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.sendPollMessage(ctx, i, poll)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно отправлен.", nil
}
