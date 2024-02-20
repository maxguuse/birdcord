package liverole

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) clearLiveRoles(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	liveroles, err := h.Database.Liveroles().GetLiveroles(ctx, guild.ID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	if len(liveroles) == 0 {
		return "Нет live-ролей.", nil
	}

	err = h.Database.Liveroles().DeleteLiveroles(
		ctx,
		guild.ID,
		lo.Map(liveroles, func(liverole *domain.Liverole, _ int) string {
			return liverole.DiscordRoleID
		}))
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Live-роли удалены.", nil
}
