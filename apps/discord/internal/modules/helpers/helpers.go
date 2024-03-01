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
	var response string
	var usersideErr *domain.UsersideError
	msg := "Произошла ошибка"

	switch {
	case errors.Is(inErr, domain.ErrInternal):
		response = "Внутренняя ошибка"
	case errors.As(inErr, &usersideErr):
		response = usersideErr.Error()
	default:
		response = "Произошла неизвестная ошибка"
	}

	_, err := session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &msg,
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Description: response,
			},
		},
	})

	return err
}
