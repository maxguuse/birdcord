package poll

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
	"github.com/maxguuse/disroute"
	"github.com/samber/lo"
)

func (h *Handler) autocompletePollList(ctx *disroute.Ctx) []*discordgo.ApplicationCommandOptionChoice {
	polls, err := h.service.GetActivePolls(ctx.Context(), &service.GetActivePollsRequest{
		GuildID: ctx.Interaction().GuildID,
		UserID:  ctx.Interaction().Member.User.ID,
	})
	if err != nil {
		h.logger.Error("Failed to get active polls", err)

		return nil
	}

	choices := lo.FilterMap(polls, func(p *domain.Poll, _ int) (*discordgo.ApplicationCommandOptionChoice, bool) {
		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("Poll ID: %d | %s", p.ID, p.Title),
			Value: p.ID,
		}

		s, ok := ctx.Options["poll"].Value.(string)
		if !ok {
			return choice, false
		}

		return choice, strings.Contains(choice.Name, s)
	})

	return choices
}

func (h *Handler) autocompleteOptionList(ctx *disroute.Ctx) []*discordgo.ApplicationCommandOptionChoice {
	rawPollId, ok := ctx.Options["poll"]
	if !ok {
		h.logger.Error("Failed to get poll id")

		return nil
	}

	poll, err := h.service.GetPoll(ctx.Context(), &service.GetPollRequest{
		GuildID: ctx.Interaction().GuildID,
		UserID:  ctx.Interaction().Member.User.ID,
		PollID:  rawPollId.IntValue(),
	})
	if err != nil {
		return nil
	}

	choices := lo.FilterMap(poll.Options,
		func(o domain.PollOption, _ int) (*discordgo.ApplicationCommandOptionChoice, bool) {
			choice := &discordgo.ApplicationCommandOptionChoice{
				Name:  o.Title,
				Value: o.ID,
			}

			s, ok := ctx.Options["option"].Value.(string)
			if !ok {
				return choice, false
			}

			return choice, strings.Contains(choice.Name, s)
		})

	return choices
}
