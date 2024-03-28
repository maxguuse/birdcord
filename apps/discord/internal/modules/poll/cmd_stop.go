package poll

import (
	"context"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) stopPoll(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	res, err := h.service.Stop(ctx, &service.StopRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  options["poll"].IntValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        res.Poll,
		interaction: i,
		stop:        true,
		fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Победители",
				Value:  strings.Join(res.Winners, ","),
				Inline: true,
			},
		},
	})
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно остановлен.", nil
}
