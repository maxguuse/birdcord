package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/samber/lo"
	"strings"
)

func (p *Polls) handleStopAutocomplete(
	s *discordgo.Session, i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	activePolls, err := p.client.GetActivePolls(context.Background(), &polls.GetActivePollsRequest{
		DiscordGuildId: i.GuildID,
	})
	if err != nil {
		fmt.Println("Error getting active polls:", err)
		return
	}
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(activePolls.Polls))
	for j, activePoll := range activePolls.Polls {
		choices[j] = &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("Poll ID: %d | %s", activePoll.Id, activePoll.Title),
			Value: activePoll.Id,
		}
	}

	err = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: lo.Filter(choices, func(c *discordgo.ApplicationCommandOptionChoice, _ int) bool {
				return strings.Contains(c.Name, options["poll"].Value.(string))
			}),
		},
	})
	if err != nil {
		fmt.Println("Error responding to activePoll:", err)
		return
	}
}
