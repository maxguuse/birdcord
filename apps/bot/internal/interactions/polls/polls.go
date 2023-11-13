package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	grpc "github.com/maxguuse/birdcord/libs/grpc/generated/polls"
)

type Polls struct {
	client      grpc.PollsClient
	commands    []*discordgo.ApplicationCommand
	subcommands []*discordgo.ApplicationCommandOption
}

var poll = &discordgo.ApplicationCommand{
	Name:        "poll",
	Description: "Управление опросами",
	Options: []*discordgo.ApplicationCommandOption{
		start,
		stop,
	},
}

func New(client grpc.PollsClient) *Polls {
	polls := &Polls{
		client: client,
		commands: []*discordgo.ApplicationCommand{
			poll,
		},
		subcommands: poll.Options,
	}

	return polls
}

func (p *Polls) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", p.commands) // bs server 1149093125118251018
	if err != nil {
		fmt.Println("Error creating poll command:", err) // Replace with logger
	}
}

func (p *Polls) CommandHandler(s *discordgo.Session, i *discordgo.Interaction) {
	commandOptions := p.buildOptionsMap(i)

	switch i.ApplicationCommandData().Options[0].Name {
	case "start":
		p.handleStartCommand(s, i, commandOptions)
	case "stop":
		p.handleStopCommand(s, i, commandOptions)
	}
}

func (p *Polls) AutocompleteHandler(s *discordgo.Session, i *discordgo.Interaction) {
	autocompleteOptions := p.buildOptionsMap(i)

	switch i.ApplicationCommandData().Options[0].Name {
	case "stop":
		p.handleStopAutocomplete(s, i, autocompleteOptions)
	}
}

func (p *Polls) buildOptionsMap(i *discordgo.Interaction) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	commandOptions := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		commandOptions[option.Name] = option
	}
	return commandOptions
}

func (p *Polls) processPollInteractionResponse(s *discordgo.Session, i *discordgo.Interaction, message *discordgo.Message, pollId int32, content string) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
	if err != nil {
		fmt.Println("Error responding to poll:", err) //TODO Replace with logger
		p.processPollFailure(s, i, message, err)
		_, invalidateErr := p.client.InvalidatePoll(context.Background(), &grpc.InvalidatePollRequest{
			PollId: pollId,
		})
		if invalidateErr != nil {
			fmt.Println("Error invalidating poll:", invalidateErr) //TODO Replace with logger
			return
		}
	}
}

func (p *Polls) processPollFailure(s *discordgo.Session, i *discordgo.Interaction, message *discordgo.Message, err error) {
	deleteErr := s.ChannelMessageDelete(message.ChannelID, message.ID)
	if deleteErr != nil {
		fmt.Println("Error deleting poll message:", deleteErr) //TODO Replace with logger
	}

	respondWithError(s, i, err)
}

func respondWithError(s *discordgo.Session, i *discordgo.Interaction, err error) {
	responseErr := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "При использовании бота возникла ошибка!",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Сообщение об ошибке",
					Description: err.Error(),
				},
			},
		},
	})
	if responseErr != nil {
		fmt.Println("Error responding to the interaction with error:", responseErr) //TODO Replace with logger
		return
	}
}
