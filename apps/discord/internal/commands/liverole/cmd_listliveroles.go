package liverole

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) listLiveRoles(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", err
	}

	liveroles, err := h.Database.Liveroles().GetLiveroles(ctx, guild.ID)
	if err != nil {
		return "", err
	}

	if len(liveroles) == 0 {
		return "Нет live-ролей.", nil
	}

	rolesList := lo.Map(liveroles, func(liverole *domain.Liverole, _ int) string {
		return fmt.Sprintf("<@&%s>", liverole.DiscordRoleID)
	})

	return "Список live-ролей: \n" + strings.Join(rolesList, "\n"), nil
}
