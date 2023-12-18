package commands

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/maxguuse/birdcord/apps/discord/internal/models"
	"github.com/maxguuse/birdcord/apps/discord/internal/postgres"
	"github.com/maxguuse/birdcord/libs/logger"
	"log/slog"
)

var poll = &discordgo.ApplicationCommand{
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
		},
	},
}

type PollCommandHandler struct {
	Log      logger.Logger
	Database *postgres.Postgres
}

func (p *PollCommandHandler) Handle(s *discordgo.Session, i interface{}) {
	cmd, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	p.Log.Info("poll command", slog.String("command", cmd.ApplicationCommandData().Name))

	commandOptions := buildCommandOptionsMap(cmd)

	switch cmd.ApplicationCommandData().Name {
	case "poll":
		p.startPoll(s, cmd, commandOptions)
	}
}

func (p *PollCommandHandler) startPoll(
	s *discordgo.Session,
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var interactionResponseContent string

	interactionResponseContent = "Опрос формируется..."
	interactionRespondErr := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: interactionResponseContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if interactionRespondErr != nil {
		p.Log.Error("error responding to interaction", slog.String("error", interactionRespondErr.Error()))
		return
	}

	transactionErr := p.Database.Transaction(func(tx pgx.Tx) error {
		getGuildQuery, getGuildQueryArgs, getGuildQueryErr := p.Database.QueryBuilder.
			Select("*").
			From("guilds").
			Where(squirrel.Eq{"discord_guild_id": i.GuildID}).
			ToSql()
		if getGuildQueryErr != nil {
			return getGuildQueryErr
		}

		var guild models.Guild

		getGuildRow := tx.QueryRow(context.Background(), getGuildQuery, getGuildQueryArgs...)
		getGuildRowErr := getGuildRow.Scan(&guild.ID, &guild.DiscordID)
		if getGuildRowErr != nil {
			return getGuildRowErr
		}

		return nil
	})
	if transactionErr != nil {
		_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Сообщение об ошибке",
					Description: transactionErr.Error(),
				},
			},
		})
		if err != nil {
			p.Log.Error("error editing interaction", slog.String("error", err.Error()))
		}
		p.Log.Error("error creating transaction", slog.String("error", transactionErr.Error()))
		return
	}

	interactionResponseContent = "Опрос создан!"
	_, interactionRespondErr = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &interactionResponseContent,
	})
	if interactionRespondErr != nil {
		p.Log.Error("error responding to interaction", slog.String("error", interactionRespondErr.Error()))
		return
	}
}
