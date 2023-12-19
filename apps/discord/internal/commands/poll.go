package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
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

	switch cmd.ApplicationCommandData().Options[0].Name {
	case "start":
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

		poll, createPolLErr := q.CreatePoll(ctx, queries.CreatePollParams{
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

		actionRows := make([]*discordgo.ActionsRow, 0, (len(optionsList)+4)/5)
		for i, option := range optionsList {
			if option == "" || utf8.RuneCountInString(option) > 50 {
				return fmt.Errorf("invalid option length")
			}

			var createOptionErr error
			pollOption, createOptionErr := q.CreatePollOption(ctx, queries.CreatePollOptionParams{
				Title:  option,
				PollID: poll.ID,
			})
			if createOptionErr != nil {
				return createOptionErr
			}

			if i%5 == 0 {
				actionRow := &discordgo.ActionsRow{
					Components: make([]discordgo.MessageComponent, 0, 5),
				}
				actionRows = append(actionRows, actionRow)
			}

			btn := discordgo.Button{
				Label:    option,
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("poll_%d_option_%d", poll.ID, pollOption.ID),
			}
			lastActionRow := actionRows[len(actionRows)-1]
			lastActionRow.Components = append(lastActionRow.Components, btn) //TODO using eventbus subscribe VoteButtonHandler here with CustomID

			optionsList[i] = fmt.Sprintf("**%d.** %s", i+1, option)
		}

		messageComponents := lo.Map(actionRows, func(actionRow *discordgo.ActionsRow, _ int) discordgo.MessageComponent {
			return actionRow
		})

		discordMsg, sendMessageErr := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       options["title"].StringValue(),
					Description: strings.Join(optionsList, "\n"),
					Timestamp:   poll.CreatedAt.Time.Format(time.RFC3339),
					Color:       0x4d58d3,
					Type:        discordgo.EmbedTypeRich,
					Author: &discordgo.MessageEmbedAuthor{
						Name:    i.Member.User.Username,
						IconURL: i.Member.User.AvatarURL(""),
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprint("Poll ID: ", poll.ID),
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Всего голосов",
							Value:  "0",
							Inline: true,
						},
					},
				},
			},
			Components: messageComponents,
		})
		if sendMessageErr != nil {
			return sendMessageErr
		}

		msg, createMessageErr := q.CreateMessage(ctx, queries.CreateMessageParams{
			DiscordMessageID: discordMsg.ID,
			DiscordChannelID: discordMsg.ChannelID,
		})
		if createMessageErr != nil {
			deleteMessageErr := s.ChannelMessageDelete(i.ChannelID, discordMsg.ID)
			return errors.Join(createMessageErr, deleteMessageErr)
		}

		createPollMessageErr := q.CreatePollMessage(ctx, queries.CreatePollMessageParams{
			PollID:    poll.ID,
			MessageID: msg.ID,
		})
		if createPollMessageErr != nil {
			deleteMessageErr := s.ChannelMessageDelete(i.ChannelID, discordMsg.ID)
			return errors.Join(createPollMessageErr, deleteMessageErr)
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
		_, interactionResponseEditErr := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
		})
		if interactionResponseEditErr != nil {
			p.Log.Error("error editing an interaction", slog.String("error", interactionResponseEditErr.Error()))
		}
	}
}
