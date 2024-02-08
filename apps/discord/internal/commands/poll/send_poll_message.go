package poll

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) sendPollMessage(
	ctx context.Context,
	i *discordgo.Interaction,
	poll *domain.PollWithDetails,
	optionsList []string,
) error {
	actionRows := h.buildActionRows(poll, i.ID)
	pollEmbed := buildPollEmbed(poll, i.Member.User)

	msg, err := h.Session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embeds:     pollEmbed,
		Components: actionRows,
	})
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	_, err = h.Database.Polls().CreatePollMessage(
		ctx, msg.ID, msg.ChannelID, poll.ID,
	)
	if err != nil {
		deleteErr := h.Session.ChannelMessageDelete(i.ChannelID, msg.ID)

		return errors.Join(domain.ErrInternal, deleteErr, err)
	}

	return nil
}

func (h *Handler) updatePollMessages(
	poll *domain.PollWithDetails,
	i *discordgo.Interaction,
	f ...*discordgo.MessageEmbedField) error {
	actionRows := h.buildActionRows(poll, i.ID)
	pollEmbed := buildPollEmbed(poll, i.Member.User)

	if len(f) > 0 {
		pollEmbed[0].Fields = append(pollEmbed[0].Fields, f...)
	}

	var wg *errgroup.Group
	for _, msg := range poll.Messages {
		msg := msg
		wg.Go(func() error {
			_, err := h.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         msg.DiscordMessageID,
				Channel:    msg.DiscordChannelID,
				Content:    new(string),
				Components: actionRows,
				Embeds:     pollEmbed,
			})

			return err
		})
	}

	if err := wg.Wait(); err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (h *Handler) buildActionRows(
	poll *domain.PollWithDetails,
	interactionID string,
) []discordgo.MessageComponent {
	buttons := make([]discordgo.MessageComponent, 0, len(poll.Options))
	for _, option := range poll.Options {
		customId := fmt.Sprintf("poll_%d_option_%d_i_%s", poll.ID, option.ID, interactionID)
		buttons = append(buttons, discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: customId,
		})

		_ = h.Pubsub.Subscribe(customId, h.BuildVoteButtonHandler(
			int32(poll.ID), int32(option.ID),
		))
	}
	buttonsGroups := lo.Chunk(buttons, 5)
	actionRows := lo.Map(buttonsGroups, func(buttons []discordgo.MessageComponent, _ int) discordgo.MessageComponent {
		return discordgo.ActionsRow{
			Components: buttons,
		}
	})

	return actionRows
}
