package liverole

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) addLiveRole(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	role := om["role"].RoleValue(h.Session, i.GuildID)

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", err
	}

	_, err = h.Database.Liveroles().CreateLiverole(ctx, role.ID, guild.ID)
	if err != nil {
		return "", err
	}

	return "Live-роль успешно добавлена.", nil
}
