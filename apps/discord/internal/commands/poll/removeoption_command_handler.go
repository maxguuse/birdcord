package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) removePollOption(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) (string, error) {
	ctx := context.Background()

	pollId := options["poll"].IntValue()
	optionId := options["option"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		return "", &domain.UsersideError{
			Msg: "Для изменения опроса нужно быть его автором.",
		}
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return "", &domain.UsersideError{
			Msg: "Опроса не существует.",
		}
	}

	err = h.Database.Polls().RemovePollOption(ctx, int(optionId))
	if err != nil {
		return "", err
	}

	poll.Options = lo.Filter(poll.Options, func(o domain.PollOption, _ int) bool {
		return o.ID != int(optionId)
	})

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", err
	}

	return "Вариант опроса успешно удален.", nil
}
