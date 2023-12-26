package helpers

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

func BuildOptionsMap(i *discordgo.Interaction) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	commandOptions := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		commandOptions[option.Name] = option
	}

	return commandOptions
}

func InteractionResponseProcess(s *discordgo.Session, i *discordgo.Interaction, msg string, err error) error {
	if err != nil {
		return interactionRespondError(err, s, i)
	}

	return interactionRespondSuccess(msg, s, i)
}

func interactionRespondSuccess(msg string, session *discordgo.Session, i *discordgo.Interaction) error {
	_, err := session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

func interactionRespondError(inErr error, session *discordgo.Session, i *discordgo.Interaction) error {
	var err error
	msg := "Произошла ошибка"

	if errors.Is(inErr, domain.ErrUserSide) {
		var response string
		switch {
		case errors.Is(inErr, domain.ErrInternal):
			response = "Внутренняя ошибка"
		case errors.Is(inErr, domain.ErrWrongPollOptionLength):
			response = "Длина варианта опроса не может быть больше 50 или меньше 1 символа"
		case errors.Is(inErr, domain.ErrAlreadyVoted):
			response = "Вы уже проголосовали в этом опросе"
		case errors.Is(inErr, domain.ErrWrongPollOptionsAmount):
			response = "Количество вариантов опроса должно быть от 2 до 25 включительно"
		case errors.Is(inErr, domain.ErrNotAuthor):
			response = "Для остановки опроса нужно быть его автором"
		case errors.Is(inErr, domain.ErrWrongGuild):
			response = "Опроса не существует"
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
