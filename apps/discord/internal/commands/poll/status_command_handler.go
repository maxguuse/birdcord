package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) statusPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) error {
	ctx := context.Background()

	pollId := options["poll"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		return errors.Join(domain.ErrUserSide, domain.ErrNotAuthor)
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return errors.Join(domain.ErrUserSide, domain.ErrWrongGuild)
	}

	return h.sendPollMessage(ctx, i, poll, lo.Map(poll.Options, func(option domain.PollOption, _ int) string {
		return option.Title
	}))
}
