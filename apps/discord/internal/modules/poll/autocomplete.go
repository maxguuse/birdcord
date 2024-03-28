package poll

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
	"github.com/samber/lo"
)

func (h *Handler) autocompletePollList(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	polls, err := h.service.GetActivePolls(ctx, &service.GetActivePollsRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
	})
	if err != nil {
		return "", err
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(polls))
	for i, poll := range polls {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("Poll ID: %d | %s", poll.ID, poll.Title),
			Value: poll.ID,
		}
	}

	err = h.Session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: lo.Filter(choices, func(c *discordgo.ApplicationCommandOptionChoice, _ int) bool {
				s, ok := options["poll"].Value.(string)
				if !ok {
					return false
				}

				return strings.Contains(c.Name, s)
			}),
		},
	})

	return "", err
}

func (h *Handler) autocompleteOptionList(i *discordgo.Interaction, options optionsMap) (string, error) {
	ctx := context.Background()

	rawPollId, ok := options["poll"]
	if !ok {
		return "", errors.New("there's no focused option")
	}

	poll, err := h.service.GetPoll(ctx, &service.GetPollRequest{
		GuildID: i.GuildID,
		UserID:  i.Member.User.ID,
		PollID:  rawPollId.IntValue(),
	})
	if err != nil {
		return "", err
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(poll.Options))
	for i, option := range poll.Options {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  option.Title,
			Value: option.ID,
		}
	}

	err = h.Session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: lo.Filter(choices, func(c *discordgo.ApplicationCommandOptionChoice, _ int) bool {
				s, ok := options["option"].Value.(string)
				if !ok {
					return false
				}

				return strings.Contains(c.Name, s)
			}),
		},
	})

	return "", err
}
