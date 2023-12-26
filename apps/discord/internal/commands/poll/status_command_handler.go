package poll

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) statusPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		if err != nil {
			h.Log.Error("error creating poll", slog.String("error", err.Error()))
			err := interactionRespondError(
				"Произошла ошибка при формировании состояния опроса",
				err, h.Session, i,
			)
			if err != nil {
				h.Log.Error(
					"error editing an interaction",
					slog.String("error", err.Error()),
				)
			}

			return
		}

		err = interactionRespondSuccess(
			"Состояние опроса сформировано!",
			h.Session, i,
		)
		if err != nil {
			h.Log.Error(
				"error editing an interaction",
				slog.String("error", err.Error()),
			)
		}
	}()

	ctx := context.Background()

	err = interactionRespondLoading(
		"Состояние опроса формируется...",
		h.Session, i,
	)
	if err != nil {
		h.Log.Error(
			"error responding to interaction",
			slog.String("error", err.Error()),
		)

		return
	}

	pollId := options["poll"].IntValue()

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

	msg, err := h.Session.ChannelMessageSend(i.ChannelID, "Bird думает...")
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	actionRows := h.buildActionRows(poll, msg, lo.Map(poll.Options, func(option domain.PollOption, _ int) string {
		return option.Title
	}))
	pollEmbed := buildPollEmbed(poll, i.Member.User)

	_, err = h.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         msg.ID,
		Channel:    msg.ChannelID,
		Content:    new(string),
		Embeds:     pollEmbed,
		Components: actionRows,
	})
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	_, err = h.Database.Polls().CreatePollMessage(
		ctx,
		msg.ID, msg.ChannelID,
		poll.ID,
	)
	if err != nil {
		deleteErr := h.Session.ChannelMessageDelete(i.ChannelID, msg.ID)
		err = errors.Join(domain.ErrInternal, deleteErr, err)

		return
	}
}