package liverole

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
)

func (h *Handler) addLiveRole(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	role := om["role"].RoleValue(h.Session, i.GuildID)

	err := h.Service.Add(ctx, &service.AddLiveRoleRequest{
		GuildID: i.GuildID,
		RoleID:  role.ID,
	})
	if err != nil {
		return "", err
	}

	return "Live-роль успешно добавлена.", nil
}
