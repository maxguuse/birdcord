package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"strings"
)

var pollStopSubcommand = &discordgo.ApplicationCommandOption{
	Name:        "stop",
	Description: "Остановить опрос",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "poll",
			Description:  "Опрос",
			Type:         discordgo.ApplicationCommandOptionInteger,
			Required:     true,
			Autocomplete: true,
		},
	},
}

func (p *Polls) handlePollStop(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	fmt.Println("Poll stopped")
	fmt.Println(i.Type)

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		pollId := options["poll"].IntValue()

		stopPollResponse, err := p.pollsClient.StopPoll(context.Background(), &polls.StopPollRequest{
			PollId: int32(pollId),
		})
		if err != nil {
			fmt.Println("Error stopping poll:", err)
			return
		}

		err = s.InteractionResponseDelete(&discordgo.Interaction{
			AppID: os.Getenv("APP_ID"),
			Token: stopPollResponse.DiscordToken,
		})
		if err != nil {
			fmt.Println("Error deleting poll:", err)
			return
		}
		pollResultsMessageData := buildPollResultsMessageData(
			stopPollResponse.Title,
			int64(stopPollResponse.TotalVotes),
			stopPollResponse.Winners,
			i.Interaction.Member,
		)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: pollResultsMessageData,
		})
		if err != nil {
			fmt.Println("Error responding to poll:", err)
			return
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		activePolls, err := p.pollsClient.GetPolls(context.Background(), &emptypb.Empty{})
		if err != nil {
			fmt.Println("Error getting active polls:", err)
			return
		}
		choices := make([]*discordgo.ApplicationCommandOptionChoice, len(activePolls.Polls))
		for i, poll := range activePolls.Polls {
			fmt.Println(poll.Title, poll.Id)
			choices[i] = &discordgo.ApplicationCommandOptionChoice{
				Name:  fmt.Sprintf("Poll ID: %d | %s", poll.Id, poll.Title),
				Value: fmt.Sprintf("%d", poll.Id),
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			fmt.Println("Error responding to poll:", err)
			return
		}
	}

}

func buildPollResultsMessageData(
	title string,
	totalVotes int64,
	options []*polls.Option,
	user *discordgo.Member,
) *discordgo.InteractionResponseData {
	var description string

	optionsTitles := lo.Map(options, func(option *polls.Option, i int) string {
		description += fmt.Sprintf("%d) %s > %d \n", i+1, option.Title, option.TotalVotes)
		return option.Title
	})
	winners := strings.Join(optionsTitles, ", ")

	embed := buildPollEmbed(title, description, user,
		nil, []*discordgo.MessageEmbedField{
			{
				Name:   "Total Votes",
				Value:  fmt.Sprintf("%d", totalVotes),
				Inline: true,
			},
			{
				Name:   "Winners",
				Value:  winners,
				Inline: true,
			},
		})

	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	}
}
