package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) startPoll(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.Create(ctx, &service.CreateRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		Poll: service.Poll{
			Title:   options["title"].StringValue(),
			Options: options["options"].StringValue(),
		},
	})
	if err != nil {
		return "", err
	}

	err = h.sendPollMessage(ctx, i, poll)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно создан.", nil
}
