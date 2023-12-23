package poll

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

func interactionRespondLoading(msg string, session *discordgo.Session, i *discordgo.Interaction) error {
	err := session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return err
}

func interactionRespondSuccess(msg string, session *discordgo.Session, i *discordgo.Interaction) error {
	_, err := session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

func interactionRespondError(msg string, inErr error, session *discordgo.Session, i *discordgo.Interaction) error {
	var err error

	if errors.Is(inErr, domain.ErrInternal) {
		_, err = session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &msg,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Description: "internal error",
				},
			},
		})
	}

	if errors.Is(inErr, domain.ErrUserSide) {
		var response string
		switch inErr {
		case domain.ErrWrongPollOptionLength:
			response = "Длина варианта опроса не может быть больше 50 символов"
		case domain.ErrAlreadyVoted:
			response = "Вы уже проголосовали в этом опросе"
		default:
			response = inErr.Error()
		}
		_, err = session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &msg,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Description: response,
				},
			},
		})
	}

	return err
}

func buildCommandOptionsMap(i *discordgo.Interaction) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	commandOptions := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		commandOptions[option.Name] = option
	}
	return commandOptions
}

func buildPollEmbed(
	poll *domain.PollWithDetails,
	optionsList []string,
	user *discordgo.User,
	votesAmount int,
) []*discordgo.MessageEmbed {
	return []*discordgo.MessageEmbed{
		{
			Title:       poll.Title,
			Description: strings.Join(optionsList, "\n"),
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
					Value:  strconv.Itoa(votesAmount),
					Inline: true,
				},
			},
		},
	}
}
