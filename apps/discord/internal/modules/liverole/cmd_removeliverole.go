package liverole

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
)

func (h *Handler) removeLiveRole(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	role := om["role"].RoleValue(h.Session, i.GuildID)

	err := h.Service.Remove(ctx, &service.RemoveLiveRoleRequest{
		GuildID: i.GuildID,
		RoleID:  role.ID,
	})
	if err != nil {
		return "", err
	}

	return "Live-роль успешно удалена.", nil
}
