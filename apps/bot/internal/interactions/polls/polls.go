package polls

import (
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
	commandOptions := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		commandOptions[option.Name] = option
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "start":
		p.handleStart(s, i, commandOptions)
	case "stop":
		//TODO implement stop poll
		fmt.Println("Not implemented: stop")
	}
}

/* Interactions used by polls
 *
 * [+] - Application command
 * Application command autocomplete
 * Message component
 */
