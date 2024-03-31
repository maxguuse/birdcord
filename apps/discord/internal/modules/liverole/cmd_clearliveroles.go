package liverole

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) clearLiveRoles(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	err := h.Service.Clear(ctx, i.GuildID)
	if err != nil {
		return "", err
	}

	return "Live-роли удалены.", nil
}
