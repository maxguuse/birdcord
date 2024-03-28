package poll

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) sendPollMessage(
	ctx context.Context,
	i *discordgo.Interaction,
	poll *domain.PollWithDetails,
) error {
	actionRows := h.buildActionRows(poll, i.ID)
	pollEmbed := buildPollEmbed(poll, i.Member.User)

	msg, err := h.Session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embeds:     pollEmbed,
		Components: actionRows,
	})
	if err != nil {
		return err
	}

	err = h.service.CreateMessage(ctx, &service.CreateMessageRequest{
		PollID: poll.ID,
		Message: service.Message{
			ID:        msg.ID,
			ChannelID: i.ChannelID,
		},
	})

	if err != nil {
		return errors.Join(h.Session.ChannelMessageDelete(i.ChannelID, msg.ID), err)
	}

	return nil
}

type UpdatePollMessageData struct {
	poll        *domain.PollWithDetails
	interaction *discordgo.Interaction
	stop        bool
	fields      []*discordgo.MessageEmbedField
}

func (h *Handler) updatePollMessages(data *UpdatePollMessageData) error {
	actionRows := lo.
		If(data.stop, make([]discordgo.MessageComponent, 0)).
		Else(h.buildActionRows(data.poll, data.interaction.ID))

	author, err := h.Session.User(data.poll.Author.DiscordUserID)
	if err != nil {
		return err
	}

	pollEmbed := buildPollEmbed(data.poll, author)

	if len(data.fields) > 0 {
		pollEmbed[0].Fields = append(pollEmbed[0].Fields, data.fields...)
	}

	wg := new(errgroup.Group)
	for _, msg := range data.poll.Messages {
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
		return err
	}

	return nil
}

func (h *Handler) buildActionRows(
	poll *domain.PollWithDetails,
	interactionID string,
) []discordgo.MessageComponent {
	buttons := make([]discordgo.MessageComponent, 0, len(poll.Options))
	for _, option := range poll.Options {
		customId := fmt.Sprintf("poll-vote-btn:poll_%d_option_%d_i_%s", poll.ID, option.ID, interactionID)
		buttons = append(buttons, discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: customId,
		})
	}
	buttonsGroups := lo.Chunk(buttons, 5)
	actionRows := lo.Map(buttonsGroups, func(buttons []discordgo.MessageComponent, _ int) discordgo.MessageComponent {
		return discordgo.ActionsRow{
			Components: buttons,
		}
	})

	return actionRows
}
