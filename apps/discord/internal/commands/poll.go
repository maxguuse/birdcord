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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
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
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

func newPolls(
	log logger.Logger,
	eb *eventbus.EventBus,
	db *db.DB,
	s *discordgo.Session,
) *PollCommandHandler {
	return &PollCommandHandler{
		Log:      log,
		Database: db,
		EventBus: eb,
		Session:  s,
	}
}

func (p *PollCommandHandler) Handle(i any) {
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

func (p *PollCommandHandler) startPoll(
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

type VoteButtonHandler struct {
	poll_id   int32
	option_id int32

	Log      logger.Logger
	Database *db.DB
	Session  *discordgo.Session
}

func (v *VoteButtonHandler) Handle(i any) {
	var err error

	vote, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	ctx := context.Background()

	interactionResponseContent := "Голос регистрируется..."
	interactionRespondErr := v.Session.InteractionRespond(vote, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: interactionResponseContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if interactionRespondErr != nil {
		v.Log.Error("error responding to interaction", slog.String("error", interactionRespondErr.Error()))
		return
	}

	tErr := v.Database.Transaction(func(q *queries.Queries) error {
		user, err := q.GetUserByDiscordID(ctx, vote.Member.User.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if user.ID == 0 {
			user, err = q.CreateUser(ctx, vote.Member.User.ID)
			if err != nil {
				return err
			}
		}

		poll, err := q.GetPoll(ctx, v.poll_id)
		if err != nil {
			return err
		}

		pollMessages, err := q.GetMessagesForPollById(ctx, poll.ID)
		if err != nil {
			return err
		}

		pollOptions, err := q.GetPollOptions(ctx, poll.ID)
		if err != nil {
			return err
		}

		err = q.AddVote(ctx, queries.AddVoteParams{
			UserID:   user.ID,
			PollID:   poll.ID,
			OptionID: v.option_id,
		})

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("you already voted for this poll")
		}
		if err != nil && !errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return err
		}

		pollVotes, err := q.GetAllVotesForPollById(ctx, poll.ID)
		if err != nil {
			return err
		}

		author, err := q.GetUserById(ctx, poll.AuthorID.Int32)
		if err != nil {
			return err
		}

		discordAuthor, err := v.Session.User(author.DiscordUserID)
		if err != nil {
			return err
		}

		optionsList := lo.Map(pollOptions, func(option queries.PollOption, i int) string {
			return fmt.Sprintf("%d. %s", i, option.Title)
		})

		for _, msg := range pollMessages {
			discordMsg, err := q.GetMessageById(ctx, msg.MessageID)
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}

			if err != nil {
				return err
			}

			pollEmbed := buildPollEmbed(
				poll,
				optionsList,
				discordAuthor,
				len(pollVotes),
			)

			_, err = s.ChannelMessageEditEmbeds(
				discordMsg.DiscordChannelID,
				discordMsg.DiscordMessageID,
				pollEmbed,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	var pgErr *pgconn.PgError

	if tErr != nil && errors.Is(tErr, pgErr) {
		v.Log.Error("error registering vote", slog.String("error", tErr.Error()))
		interactionResponseContent = "Произошла внутренняя ошибка при регистрации голоса"
		_, err := v.Session.InteractionResponseEdit(vote, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
		})
		if err != nil {
			v.Log.Error("error editing an interaction", slog.String("error", err.Error()))
		}

		return
	}

	if tErr != nil {
		v.Log.Info("error registering vote", slog.String("error", tErr.Error()))
		interactionResponseContent = "Произошла ошибка при регистрации голоса"
		_, err := v.Session.InteractionResponseEdit(vote, &discordgo.WebhookEdit{
			Content: &interactionResponseContent,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Description: tErr.Error(),
				},
			},
		})
		if err != nil {
			v.Log.Error("error editing an interaction", slog.String("error", err.Error()))
		}

		return
	}

	interactionResponseContent = "Голос засчитан!"
	_, err = v.Session.InteractionResponseEdit(vote, &discordgo.WebhookEdit{
		Content: &interactionResponseContent,
	})
	if err != nil {
		v.Log.Error("error editing an interaction", slog.String("error", err.Error()))
	}
}
