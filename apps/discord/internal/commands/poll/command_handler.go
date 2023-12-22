package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
)

type CommandHandler struct {
	Log      logger.Logger
	Database *db.DB
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

func NewCommandHandler(
	log logger.Logger,
	eb *eventbus.EventBus,
	db *db.DB,
	s *discordgo.Session,
) *CommandHandler {
	return &CommandHandler{
		Log:      log,
		Database: db,
		EventBus: eb,
		Session:  s,
	}
}

func (p *CommandHandler) Handle(i any) {
	cmd, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	commandOptions := buildCommandOptionsMap(cmd)

	switch cmd.ApplicationCommandData().Options[0].Name {
	case "start":
		p.startPoll(cmd, commandOptions)
	}
}

func (p *CommandHandler) startPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	ctx := context.Background()

	interactionResponseContent := "Опрос формируется..."
	interactionRespondErr := p.Session.InteractionRespond(i, &discordgo.InteractionResponse{
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

	err := p.Database.Transaction(func(q *queries.Queries) error {
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
			return fmt.Errorf("not enough options")
		}

		buttons := make([]discordgo.MessageComponent, 0, len(optionsList))
		for i, option := range optionsList {
			if len(option) == 0 || utf8.RuneCountInString(option) > 50 {
				return fmt.Errorf("invalid option length")
			}

			pollOption, err := q.CreatePollOption(ctx, queries.CreatePollOptionParams{
				Title:  option,
				PollID: poll.ID,
			})
			if err != nil {
				return err //TODO add ErrQueryFailed in db lib
			}

			customId := fmt.Sprintf("poll_%d_option_%d", poll.ID, pollOption.ID)
			buttons = append(buttons, discordgo.Button{
				Label:    pollOption.Title,
				Style:    discordgo.PrimaryButton,
				CustomID: customId,
			})

			p.EventBus.Subscribe(customId, &VoteButtonHandler{
				poll_id:   poll.ID,
				option_id: pollOption.ID,
				Log:       p.Log,
				Database:  p.Database,
				Session:   p.Session,
			})

			optionsList[i] = fmt.Sprintf("**%d.** %s", i+1, option)
		}
		buttonsGroups := lo.Chunk(buttons, 5)
		actionRows := lo.Map(buttonsGroups, func(buttons []discordgo.MessageComponent, _ int) discordgo.MessageComponent {
			return discordgo.ActionsRow{
				Components: buttons,
			}
		})

		a := buildPollEmbed(
			poll,
			optionsList,
			i.Member.User,
			0,
		)
		discordMsg, sendMessageErr := p.Session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Embeds:     a,
			Components: actionRows,
		})
		if sendMessageErr != nil {
			return sendMessageErr
		}

		msg, createMessageErr := q.CreateMessage(ctx, queries.CreateMessageParams{
			DiscordMessageID: discordMsg.ID,
			DiscordChannelID: discordMsg.ChannelID,
		})
		if createMessageErr != nil {
			deleteMessageErr := p.Session.ChannelMessageDelete(i.ChannelID, discordMsg.ID)
			return errors.Join(createMessageErr, deleteMessageErr)
		}

		createPollMessageErr := q.CreatePollMessage(ctx, queries.CreatePollMessageParams{
			PollID:    poll.ID,
			MessageID: msg.ID,
		})
		if createPollMessageErr != nil {
			deleteMessageErr := p.Session.ChannelMessageDelete(i.ChannelID, discordMsg.ID)
			return errors.Join(createPollMessageErr, deleteMessageErr)
		}

		return nil
	})

	var pgErr *pgconn.PgError

	if err != nil && errors.Is(err, pgErr) {
		p.Log.Error("error creating poll", slog.String("error", err.Error()))
		interactionResponseContent = "Произошла внутренняя ошибка при формировании опроса"
		_, err := p.Session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
		})
		if err != nil {
			p.Log.Error("error editing an interaction", slog.String("error", err.Error()))
		}

		return
	}

	if err != nil {
		p.Log.Error("error creating poll", slog.String("error", err.Error()))
		interactionResponseContent = "Произошла ошибка при формировании опроса"
		_, err := p.Session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Description: err.Error(),
				},
			},
		})
		if err != nil {
			p.Log.Error("error editing an interaction", slog.String("error", err.Error()))
		}

		return
	}

	interactionResponseContent = "Опрос создан!"
	_, interactionResponseEditErr := p.Session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &interactionResponseContent,
	})
	if interactionResponseEditErr != nil {
		p.Log.Error("error editing an interaction", slog.String("error", interactionResponseEditErr.Error()))
	}
}
