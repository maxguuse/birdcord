package poll

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
	"github.com/maxguuse/disroute"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) sendPollMessage(
	ctx *disroute.Ctx,
	poll *domain.PollWithDetails,
) error {
	actionRows := buildActionRows(poll, ctx.Interaction().ID)
	pollEmbed := buildPollEmbed(poll, ctx.Interaction().Member.User)

	msg, err := ctx.Session().ChannelMessageSendComplex(ctx.Interaction().ChannelID, &discordgo.MessageSend{
		Embeds:     pollEmbed,
		Components: actionRows,
	})
	if err != nil {
		return err
	}

	err = h.service.CreateMessage(ctx.Context(), &service.CreateMessageRequest{
		PollID: poll.ID,
		Message: service.Message{
			ID:        msg.ID,
			ChannelID: ctx.Interaction().ChannelID,
		},
	})

	if err != nil {
		return errors.Join(ctx.Session().ChannelMessageDelete(ctx.Interaction().ChannelID, msg.ID), err)
	}

	return nil
}

type UpdatePollMessageData struct {
	poll   *domain.PollWithDetails
	ctx    *disroute.Ctx
	stop   bool
	fields []*discordgo.MessageEmbedField
}

func (h *Handler) updatePollMessages(data *UpdatePollMessageData) error {
	actionRows := lo.
		If(data.stop, make([]discordgo.MessageComponent, 0)).
		Else(buildActionRows(data.poll, data.ctx.Interaction().ID))

	author, err := data.ctx.Session().User(strconv.Itoa(data.poll.AuthorID))
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
			_, err := data.ctx.Session().ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         msg.DiscordMessageID,
				Channel:    msg.DiscordChannelID,
				Content:    new(string),
				Components: &actionRows,
				Embeds:     &pollEmbed,
			})

			return err
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return nil
}

func buildActionRows(
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

const (
	VOTES_BAR_BLOCK = "■"
	VOTES_BAR_SPACE = " "
)

func buildPollEmbed(
	poll *domain.PollWithDetails,
	user *discordgo.User,
) []*discordgo.MessageEmbed {
	totalVotes := len(poll.Votes)

	optionsList := lo.Map(poll.Options, func(option domain.PollOption, i int) string {
		return fmt.Sprintf("**%d**. %s", i+1, option.Title)
	})

	optionsPercentageBars := lo.Map(poll.Options, func(option domain.PollOption, i int) string {
		votesForOption := lo.CountBy(poll.Votes, func(vote domain.PollVote) bool {
			return vote.OptionID == option.ID
		})

		percentage := (float64(votesForOption) / float64(totalVotes)) * 100
		if math.IsNaN(percentage) {
			percentage = 0
		}

		t := math.Ceil(percentage)
		t2 := int(math.Floor(t / 3.33))

		if t2 < 0 {
			t2 = 0
		}

		bar := strings.Repeat(VOTES_BAR_BLOCK, t2) + strings.Repeat(VOTES_BAR_SPACE, 30-t2)

		return fmt.Sprintf("(%d) | %s | (%d%%)", i+1, bar, int(t))
	})

	optionsListDesc := strings.Join(optionsList, "\n")
	optionsBarsDesc := strings.Join(optionsPercentageBars, "\n")

	return []*discordgo.MessageEmbed{
		{
			Title:       poll.Title,
			Description: optionsListDesc + "\n```" + optionsBarsDesc + "```",
			Timestamp:   poll.CreatedAt.Format(time.RFC3339),
			Color:       0x4d58d3,
			Type:        discordgo.EmbedTypeRich,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.Username,
				IconURL: user.AvatarURL(""),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprint("Poll ID: ", poll.ID),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Всего голосов",
					Value:  strconv.Itoa(totalVotes),
					Inline: true,
				},
			},
		},
	}
}
