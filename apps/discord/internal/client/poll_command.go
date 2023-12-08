package client

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

type PollCommand struct{}

var (
	poll = &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Управление опросами",
		Options: []*discordgo.ApplicationCommandOption{
			start,
			stop,
		},
	}

	start = &discordgo.ApplicationCommandOption{
		Name:        "start",
		Description: "Начать опрос",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "Заголовок опроса",
				Type:        discordgo.ApplicationCommandOptionString,
				MaxLength:   50,
				Required:    true,
			},
			{
				Name:        "options",
				Description: "Варианты ответа (разделите их символом '|')",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}

	stop = &discordgo.ApplicationCommandOption{
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
)

func (p *PollCommand) GetCommand() *discordgo.ApplicationCommand {
	return poll
}

func (p *PollCommand) Execute(c *Client, i *discordgo.InteractionCreate) {
	c.Log.Debug(
		"Got poll command",
		slog.String("command", i.ApplicationCommandData().Name),
	)
}
