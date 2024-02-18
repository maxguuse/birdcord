package poll

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (h *Handler) autocompletePollList(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	ctx := context.Background()

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return
	}

	polls, err := h.Database.Polls().GetActivePolls(ctx, guild.ID, user.ID)
	if err != nil {
		return
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
	if err != nil {
		return
	}
}

func (h *Handler) autocompleteOptionList(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	ctx := context.Background()

	pollId := options["poll"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		return
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
	if err != nil {
		return
	}
}
