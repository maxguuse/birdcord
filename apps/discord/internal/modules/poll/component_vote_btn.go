package poll

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) VoteBtnHandler(i *discordgo.Interaction) (string, error) {
	ctx := context.Background()

	poll, err := h.service.AddVote(ctx, &service.AddVoteRequest{
		UserID:   i.Member.User.ID,
		CustomID: i.MessageComponentData().CustomID,
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})

	return "Голос зарегистрирован.", err
}
