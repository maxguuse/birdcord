package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
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
	Database *db.DB
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
	ctx := context.Background()
	var interactionResponseContent string

	interactionResponseContent = "Опрос формируется..."
	createMessageErr := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: interactionResponseContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if createMessageErr != nil {
		p.Log.Error("error responding to interaction", slog.String("error", createMessageErr.Error()))
		return
	}

	transactionErr, callbackErr := p.Database.Transaction(func(q *queries.Queries) error {
		guild, getGuildErr := q.GetGuildByDiscordID(ctx, i.GuildID)
		if getGuildErr != nil {
			return getGuildErr
		}

		user, getUserErr := q.GetUserByDiscordID(ctx, i.Member.User.ID)
		if getUserErr != nil && !errors.Is(getUserErr, pgx.ErrNoRows) {
			return getUserErr
		}
		if user.ID == 0 {
			var createUserErr error
			user, createUserErr = q.CreateUser(ctx, i.Member.User.ID)
			if createUserErr != nil {
				return createUserErr
			}
		}

		pollId, createPolLErr := q.CreatePoll(ctx, queries.CreatePollParams{
			Title: options["title"].StringValue(),
			AuthorID: pgtype.Int4{
				Int32: user.ID,
				Valid: true,
			},
			GuildID: guild.ID,
		})
		if createPolLErr != nil {
			return createPolLErr
		}

		rawOptions := options["options"].StringValue()
		optionsList := strings.Split(rawOptions, "|")
		if len(optionsList) < 2 {
			return nil
		}

		for _, option := range optionsList {
			if option == "" || utf8.RuneCountInString(option) > 50 {
				return fmt.Errorf("invalid option length")
			}
			_, createOptionErr := q.CreatePollOption(ctx, queries.CreatePollOptionParams{
				Title:  option,
				PollID: pollId,
			})
			if createOptionErr != nil {
				return createOptionErr
			}
		}

		discordMsg, sendMessageErr := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Content: "OK!",
		})
		if createMessageErr != nil {
			return sendMessageErr
		}

		msg, createMessageErr := q.CreateMessage(ctx, queries.CreateMessageParams{
			DiscordMessageID: discordMsg.ID,
			DiscordChannelID: discordMsg.ChannelID,
		})
		if createMessageErr != nil {
			return createMessageErr
		}

		createPollMessageErr := q.CreatePollMessage(ctx, queries.CreatePollMessageParams{
			PollID:    pollId,
			MessageID: msg.ID,
		})
		if createPollMessageErr != nil {
			return createPollMessageErr
		}

		return nil
	})

	if transactionErr != nil && !errors.Is(transactionErr, pgx.ErrTxClosed) {
		p.Log.Error("error creating poll", slog.String("error", transactionErr.Error()))
		interactionResponseContent = "Произошла внутренняя ошибка при формировании опроса"
		_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
		})
		if err != nil {
			p.Log.Error("error editing an interaction", slog.String("error", err.Error()))
		}
	} else if callbackErr != nil {
		p.Log.Error("error creating poll", slog.String("error", callbackErr.Error()))
		interactionResponseContent = "Произошла ошибка при формировании опроса"
		_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Description: callbackErr.Error(),
				},
			},
		})
		if err != nil {
			p.Log.Error("error editing an interaction", slog.String("error", err.Error()))
		}
	} else {
		interactionResponseContent = "Опрос создан!"
		_, createMessageErr = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
		})
		if createMessageErr != nil {
			p.Log.Error("error editing an interaction", slog.String("error", createMessageErr.Error()))
		}
	}
}
