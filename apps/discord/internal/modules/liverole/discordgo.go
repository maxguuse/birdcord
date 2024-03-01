package liverole

import "github.com/bwmarrin/discordgo"

func (h *Handler) GetDiscordGo() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		command,
	}
}

var command = &discordgo.ApplicationCommand{
	Name:        "liverole",
	Description: "Управление live-ролями",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        SubcommandAdd,
			Description: "Добавить live-роль",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "role",
					Description:  "Роль",
					Type:         discordgo.ApplicationCommandOptionRole,
					Autocomplete: true,
					Required:     true,
				},
			},
		},
		{
			Name:        SubcommandRemove,
			Description: "Удалить live-роль",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "role",
					Description:  "Роль",
					Type:         discordgo.ApplicationCommandOptionRole,
					Autocomplete: true,
					Required:     true,
				},
			},
		},
		{
			Name:        SubcommandList,
			Description: "Список live-ролей",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        SubcommandClear,
			Description: "Очистить список live-ролей",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
	},
}
