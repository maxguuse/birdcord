package scommands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/types"
	"strings"
)

var polls = []*discordgo.ApplicationCommand{
	{
		Name:        "poll",
		Description: "Управление опросами",
		Options: []*discordgo.ApplicationCommandOption{
			{
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
			},
		},
	},
}

func PollCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := buildOptionsMap(i.ApplicationCommandData().Options[0].Options)

	poll := &types.Poll{
		ID:      0,
		Title:   options["title"].StringValue(),
		Options: strings.Split(options["options"].StringValue(), "|"),
	}

	var pollMessageData *discordgo.InteractionResponseData
	if len(poll.Options) > 25 {
		pollMessageData = &discordgo.InteractionResponseData{
			Content: "Слишком много опций!",
		}
	} else {
		pollMessageData = buildPollMessage(poll.Title, poll.Options, i.Member.User)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: pollMessageData,
	})
	if err != nil {
		fmt.Println("Error responding to poll:", err)
	}
}

func registerPollsCommands(s *discordgo.Session) {
	for _, command := range polls {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			fmt.Println("Error creating poll command:", err)
		}
	}
}
