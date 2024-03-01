package liverole

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
)

func (h *Handler) addLiveRole(
	i *discordgo.Interaction,
	om optionsMap,
) (string, error) {
	ctx := context.Background()

	role := om["role"].RoleValue(h.Session, i.GuildID)

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	_, err = h.Database.Liveroles().CreateLiverole(ctx, role.ID, guild.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return "", &domain.UsersideError{
				Msg: "Данная роль уже добавлена.",
			}
		} else {
			return "", errors.Join(domain.ErrInternal, err)
		}
	}

	return "Live-роль успешно добавлена.", nil
}
