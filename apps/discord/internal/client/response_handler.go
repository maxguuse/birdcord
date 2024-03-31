package client

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/disroute"
)

func (c *Client) responseHandler(ctx *disroute.Ctx, r *disroute.Response) {
	if r.Err != nil {
		err := c.interactionRespondError(r.Err, ctx.Session(), ctx.Interaction())
		if err != nil {
			c.logger.Error("error responding to interaction", slog.Any("error", err))
		}

		return
	}

	if r.CustomResponse != nil {
		_, err := ctx.Session().InteractionResponseEdit(ctx.Interaction(), responseToWebhook(r.CustomResponse))
		if err != nil {
			c.logger.Error("error responding to interaction", slog.Any("error", err))
		}

		return
	}

	_, err := ctx.Session().InteractionResponseEdit(ctx.Interaction(), &discordgo.WebhookEdit{
		Content: &r.Message,
	})
	if err != nil {
		c.logger.Error("error editting response to interaction", slog.Any("error", err))
	}
}

func (c *Client) interactionRespondError(inErr error, session *discordgo.Session, i *discordgo.Interaction) error {
	var response string
	var usersideErr *domain.UsersideError
	msg := "Произошла ошибка"

	switch {
	case errors.Is(inErr, domain.ErrInternal):
		c.logger.Error("internal error", slog.Any("error", inErr))
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

func responseToWebhook(r *discordgo.InteractionResponse) *discordgo.WebhookEdit {
	return &discordgo.WebhookEdit{
		Content:         &r.Data.Content,
		Components:      &r.Data.Components,
		Embeds:          &r.Data.Embeds,
		Files:           r.Data.Files,
		AllowedMentions: r.Data.AllowedMentions,
	}
}
