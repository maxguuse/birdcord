package poll

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
	"github.com/maxguuse/disroute"
)

func (h *Handler) addOption(ctx *disroute.Ctx) disroute.Response {
	poll, err := h.service.AddOption(ctx.Context(), &service.AddOptionRequest{
		GuildID: ctx.Interaction().GuildID,
		UserID:  ctx.Interaction().Member.User.ID,
		PollID:  ctx.Options["poll"].IntValue(),
		Option:  ctx.Options["option"].StringValue(),
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll: poll,
		ctx:  ctx,
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	return disroute.Response{
		Message: "Вариант опроса успешно добавлен.",
	}
}

func (h *Handler) removeOption(ctx *disroute.Ctx) disroute.Response {
	poll, err := h.service.RemoveOption(ctx.Context(), &service.RemoveOptionRequest{
		GuildID:  ctx.Interaction().GuildID,
		UserID:   ctx.Interaction().Member.User.ID,
		PollID:   ctx.Options["poll"].IntValue(),
		OptionID: ctx.Options["option"].IntValue(),
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll: poll,
		ctx:  ctx,
	})
	if err != nil {
		return disroute.Response{
			Err: errors.Join(domain.ErrInternal, err),
		}
	}

	return disroute.Response{
		Message: "Вариант опроса успешно удалён.",
	}
}

func (h *Handler) start(ctx *disroute.Ctx) disroute.Response {
	poll, err := h.service.Create(ctx.Context(), &service.CreateRequest{
		GuildID: ctx.Interaction().GuildID,
		UserID:  ctx.Interaction().Member.User.ID,
		Poll: service.Poll{
			Title:   ctx.Options["title"].StringValue(),
			Options: ctx.Options["options"].StringValue(),
		},
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	err = h.sendPollMessage(ctx, poll)
	if err != nil {
		return disroute.Response{
			Err: errors.Join(domain.ErrInternal, err),
		}
	}

	return disroute.Response{
		Message: "Опрос успешно создан.",
	}
}

func (h *Handler) status(ctx *disroute.Ctx) disroute.Response {
	poll, err := h.service.GetPoll(ctx.Context(), &service.GetPollRequest{
		GuildID: ctx.Interaction().GuildID,
		UserID:  ctx.Interaction().Member.User.ID,
		PollID:  ctx.Options["poll"].IntValue(),
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	err = h.sendPollMessage(ctx, poll)
	if err != nil {
		return disroute.Response{
			Err: errors.Join(domain.ErrInternal, err),
		}
	}

	return disroute.Response{
		Message: "Опрос успешно отправлен.",
	}
}

func (h *Handler) stop(ctx *disroute.Ctx) disroute.Response {
	res, err := h.service.Stop(ctx.Context(), &service.StopRequest{
		GuildID: ctx.Interaction().GuildID,
		UserID:  ctx.Interaction().Member.User.ID,
		PollID:  ctx.Options["poll"].IntValue(),
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll: res.Poll,
		ctx:  ctx,
		stop: true,
		fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Победители",
				Value:  strings.Join(res.Winners, ","),
				Inline: true,
			},
		},
	})
	if err != nil {
		return disroute.Response{
			Err: errors.Join(domain.ErrInternal, err),
		}
	}

	return disroute.Response{
		Message: "Опрос успешно остановлен.",
	}
}
