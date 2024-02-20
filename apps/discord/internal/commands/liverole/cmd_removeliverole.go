package liverole

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

func (h *Handler) removeLiveRole(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	role := om["role"].RoleValue(h.Session, i.GuildID)

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	r, err := h.Database.Liveroles().GetLiverole(ctx, guild.ID, role.ID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	err = h.Database.Liveroles().DeleteLiverole(ctx, r.ID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Live-роль успешно удалена.", nil
}
