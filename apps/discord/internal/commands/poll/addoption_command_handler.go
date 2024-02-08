package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

func (h *Handler) addPollOption(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) (string, error) {
	ctx := context.Background()

	pollId := options["poll"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		return "", errors.Join(domain.ErrUserSide, domain.ErrNotAuthor)
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return "", errors.Join(domain.ErrUserSide, domain.ErrWrongGuild)
	}

	newOption, err := h.Database.Polls().AddPollOption(ctx, int(pollId), options["option"].StringValue())
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	poll.Options = append(poll.Options, *newOption)

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", err
	}

	return "Вариант опроса успешно добавлен.", nil
}
