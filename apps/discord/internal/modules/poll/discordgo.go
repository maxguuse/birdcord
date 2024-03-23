package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (h *Handler) GetDiscordGo() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		command,
	}
}

var command = &discordgo.ApplicationCommand{
	Name:         CommandPoll,
	Description:  "Управление опросами",
	DMPermission: lo.ToPtr(false),
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        SubcommandStart,
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
		},
		{
			Name:        SubcommandStop,
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
		},
		{
			Name:        SubcommandStatus,
			Description: "Статус опроса",
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
		},
		{
			Name:        SubcommandAddOption,
			Description: "Добавить вариант ответа к опросу",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "poll",
					Description:  "Опрос",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:        "option",
					Description: "Новый вариант ответа",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					MaxLength:   50,
				},
			},
		},
		{
			Name:        SubcommandRemoveOption,
			Description: "Удалить вариант ответа из опроса",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "poll",
					Description:  "Опрос",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "option",
					Description:  "Вариант ответа",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	},
}
