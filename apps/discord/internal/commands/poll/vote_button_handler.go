package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
)

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

			_, err = v.Session.ChannelMessageEditEmbeds(
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
