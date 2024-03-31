package liverole

import (
	"context"
	"strings"

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
