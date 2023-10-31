package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"strings"
)

var pollStartSubcommand = &discordgo.ApplicationCommandOption{
	Name:        "start",
	Description: "Начать опрос",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "title",
			Description: "Заголовок опроса",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		}, // Title of the poll
		{
			Name:        "options",
			Description: "Варианты ответа (разделите их символом '|')",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func (p *Polls) handlePollStart(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	pollOptions := strings.Split(options["options"].StringValue(), "|")
	if len(pollOptions) > 25 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Слишком много опций!",
			},
		})
		if err != nil {
			fmt.Println("Error responding to poll:", err)
		}
		return
	}

	createPollResponse, err := p.pollsClient.CreatePoll(context.Background(), &polls.CreatePollRequest{
		Title:        options["title"].StringValue(),
		Options:      pollOptions,
		DiscordToken: i.Interaction.Token,
	})
	if err != nil {
		fmt.Println("Error creating poll:", err)
		return
	}

	pollMessageData := buildPollMessageData(
		options["title"].StringValue(),
		fmt.Sprintf("%d", createPollResponse.PollId),
		0,
		createPollResponse.GetOptions(),
		i.Interaction.Member,
	)

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: pollMessageData,
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println("Error responding to poll:", err)
		return
	}
}

func buildPollMessageData(
	title string,
	pollId string,
	totalVotes int32,
	options []*polls.Option,
	user *discordgo.Member,
) *discordgo.InteractionResponseData {
	var description string
	buttons := make([]discordgo.MessageComponent, len(options))
	for i, option := range options {
		description += fmt.Sprintf("%d) %s > %d \n", i+1, option.Title, option.TotalVotes)
		buttons[i] = discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("poll_%s_choice_%d", pollId, option.Id),
		}
	}

	actionRows := make([]discordgo.MessageComponent, (len(buttons)+4)/5)
	for i := 0; i < len(buttons); i += 5 {
		actionRow := discordgo.ActionsRow{}
		for j := 0; j < 5; j++ {
			if i+j >= len(buttons) {
				break
			}
			actionRow.Components = append(actionRow.Components, buttons[i+j])
		}
		actionRows[i/5] = actionRow
	}

	embed := buildPollEmbed(title, description, user,
		&discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Poll ID: %s", pollId),
		},
		[]*discordgo.MessageEmbedField{
			{
				Name:   "Total Votes",
				Value:  fmt.Sprintf("%d", totalVotes),
				Inline: true,
			},
		},
	)

	responseData := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
		Components: actionRows,
	}
	return responseData
}
