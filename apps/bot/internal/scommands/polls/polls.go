package polls

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	grpcPolls "github.com/maxguuse/birdcord/libs/grpc/generated/polls"
)

type Polls struct {
	pollsClient grpcPolls.PollsClient

	commands []*discordgo.ApplicationCommand
}

var pollCommand = &discordgo.ApplicationCommand{
	Name:        "poll",
	Description: "Управление опросами",
	Options:     []*discordgo.ApplicationCommandOption{},
}

func New(pollsClient grpcPolls.PollsClient) *Polls {
	return &Polls{
		pollsClient: pollsClient,
		commands: []*discordgo.ApplicationCommand{
			pollCommand,
		},
	}
}

func (p *Polls) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", p.commands) // bs server 1149093125118251018
	if err != nil {
		fmt.Println("Error creating poll command:", err)
	}
}

func (p *Polls) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		p.buttonHandler(s, i)
	}
}

func (p *Polls) buttonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vote := parseVoteFromButtonInteraction(i)

	p.handleVoteButton(s, i, vote)
}
