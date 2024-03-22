package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

func (h *Handler) startPoll(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	optionsList, err := processPollOptions(options["options"].StringValue())
	if err != nil {
		return "", err
	}

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	poll, err := h.Database.Polls().CreatePoll(
		ctx,
		options["title"].StringValue(),
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	err = h.sendPollMessage(ctx, i, poll)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно создан.", nil
}
