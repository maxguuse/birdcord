package poll

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
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

		_ = h.Pubsub.Subscribe(customId, h.VoteBuilder.Build(
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
