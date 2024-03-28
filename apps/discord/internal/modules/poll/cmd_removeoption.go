package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) removePollOption(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.RemoveOption(ctx, &service.RemoveOptionRequest{
		GuildID:  i.GuildID,
		UserID:   i.Member.User.ID,
		PollID:   options["poll"].IntValue(),
		OptionID: options["option"].IntValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Вариант опроса успешно удален.", nil
}
