package poll

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) addPollOption(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.AddOption(ctx, &service.AddOptionRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  options["poll"].IntValue(),
		Option:  options["option"].StringValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", err
	}

	return "Вариант опроса успешно добавлен.", nil
}
