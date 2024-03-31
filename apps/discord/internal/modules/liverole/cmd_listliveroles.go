package liverole

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) listLiveRoles(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	rolesList, err := h.Service.List(ctx, i.GuildID)
	if err != nil {
		return "", err
	}

	return "Список live-ролей: \n" + strings.Join(rolesList, "\n"), nil
}
