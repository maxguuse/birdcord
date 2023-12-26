package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) startPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		if err != nil {
			h.Log.Error("error creating poll", slog.String("error", err.Error()))
			err := interactionRespondError(
				"Произошла ошибка при создании опроса",
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
			"Опрос создан!",
			h.Session, i,
		)
		if err != nil {
			h.Log.Error(
				"error editing an interaction",
				slog.String("error", err.Error()),
			)
		}
	}()

	err = interactionRespondLoading(
		"Опрос формируется...",
		h.Session, i,
	)
	if err != nil {
		h.Log.Error(
			"error responding to interaction",
			slog.String("error", err.Error()),
		)

		return
	}

	ctx := context.Background()

	optionsList, err := processPollOptions(options["options"].StringValue())
	if err != nil {
		return
	}

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return
	}

	poll, err := h.Database.Polls().CreatePoll(
		ctx,
		options["title"].StringValue(),
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return
	}

	msg, err := h.Session.ChannelMessageSend(i.ChannelID, "Bird думает...")
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	actionRows := h.buildActionRows(poll, msg, optionsList)
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

func (h *Handler) buildActionRows(
	poll *domain.PollWithDetails,
	msg *discordgo.Message,
	optionsList []string,
) []discordgo.MessageComponent {
	buttons := make([]discordgo.MessageComponent, 0, len(poll.Options))
	for i, option := range poll.Options {
		customId := fmt.Sprintf("poll_%d_option_%d_msg_%s", poll.ID, option.ID, msg.ID)
		buttons = append(buttons, discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: customId,
		})

		_ = h.Pubsub.Subscribe(customId, h.VoteBuilder.Build(
			int32(poll.ID), int32(option.ID),
		))

		optionsList[i] = fmt.Sprintf("**%d.** %s", i+1, option.Title)
	}
	buttonsGroups := lo.Chunk(buttons, 5)
	actionRows := lo.Map(buttonsGroups, func(buttons []discordgo.MessageComponent, _ int) discordgo.MessageComponent {
		return discordgo.ActionsRow{
			Components: buttons,
		}
	})

	return actionRows
}

func processPollOptions(rawOptions string) ([]string, error) {
	optionsList := strings.Split(rawOptions, "|")
	if len(optionsList) < 2 || len(optionsList) > 25 {
		return nil, errors.Join(
			domain.ErrUserSide,
			domain.ErrWrongPollOptionsAmount,
		)
	}
	if lo.SomeBy(optionsList, func(o string) bool {
		return utf8.RuneCountInString(o) > 50 || utf8.RuneCountInString(o) < 1
	}) {
		return nil, errors.Join(
			domain.ErrUserSide,
			domain.ErrWrongPollOptionLength,
		)
	}

	return optionsList, nil
}
