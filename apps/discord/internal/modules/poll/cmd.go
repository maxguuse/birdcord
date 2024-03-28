package poll

import (
	"context"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
)

func (h *Handler) addOption(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.AddOption(ctx, &service.AddOptionRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  options["poll"].IntValue(),
		Option:  options["option"].StringValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", err
	}

	return "Вариант опроса успешно добавлен.", nil
}

func (h *Handler) removeOption(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.RemoveOption(ctx, &service.RemoveOptionRequest{
		GuildID:  i.GuildID,
		UserID:   i.Member.User.ID,
		PollID:   options["poll"].IntValue(),
		OptionID: options["option"].IntValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Вариант опроса успешно удален.", nil
}

func (h *Handler) start(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.Create(ctx, &service.CreateRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		Poll: service.Poll{
			Title:   options["title"].StringValue(),
			Options: options["options"].StringValue(),
		},
	})
	if err != nil {
		return "", err
	}

	err = h.sendPollMessage(ctx, i, poll)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно создан.", nil
}

func (h *Handler) status(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	poll, err := h.service.GetPoll(ctx, &service.GetPollRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  options["poll"].IntValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.sendPollMessage(ctx, i, poll)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно отправлен.", nil
}

func (h *Handler) stop(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	res, err := h.service.Stop(ctx, &service.StopRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  options["poll"].IntValue(),
	})
	if err != nil {
		return "", err
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        res.Poll,
		interaction: i,
		stop:        true,
		fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Победители",
				Value:  strings.Join(res.Winners, ","),
				Inline: true,
			},
		},
	})
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return "Опрос успешно остановлен.", nil
}
