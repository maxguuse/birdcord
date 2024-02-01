package poll

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) removePollOption(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		err = helpers.InteractionResponseProcess(h.Session, i, "Вариант удален.", err)
		if err != nil {
			h.Log.Error("error editing an interaction response", slog.String("error", err.Error()))
		}
	}()

	ctx := context.Background()

	pollId := options["poll"].IntValue()
	optionId := options["option"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		err = errors.Join(domain.ErrUserSide, domain.ErrNotAuthor)

		return
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		err = errors.Join(domain.ErrUserSide, domain.ErrWrongGuild)

		return
	}

	err = h.Database.Polls().RemovePollOption(ctx, int(optionId))
	if err != nil {
		return
	}

	poll.Options = lo.Filter(poll.Options, func(o domain.PollOption, _ int) bool {
		return o.ID != int(optionId)
	})

	pollEmbed := buildPollEmbed(poll, i.Member.User)
	actionRows := h.buildActionRows(poll, i.ID)

	for _, msg := range poll.Messages {
		_, err := h.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         msg.DiscordMessageID,
			Channel:    msg.DiscordChannelID,
			Content:    new(string),
			Components: actionRows,
			Embeds:     pollEmbed,
		})
		if err != nil {
			err = errors.Join(domain.ErrInternal, err)
			h.Log.Error("error editing poll message", slog.String("error", err.Error()))
		}
	}
}
