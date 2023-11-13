package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
	"strings"
	"unicode/utf8"
)

func (p PollsServer) createPoll(ctx context.Context, request *polls.CreatePollRequest) (*polls.CreatePollResponse, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return &polls.CreatePollResponse{}, err
	}
	qr := queries.New(tx)

	pollOptions := strings.Split(request.Options, "|")
	if len(pollOptions) < 2 || len(pollOptions) > 25 {
		return &polls.CreatePollResponse{}, errors.New("количество вариантов ответа должно быть от 2 до 25")
	}

	countEmptyOptions := lo.CountBy(pollOptions, func(option string) bool {
		return option == ""
	})
	if countEmptyOptions > 0 {
		return &polls.CreatePollResponse{}, errors.New("варианты ответа не могут быть пустыми")
	}

	// Create new record in polls table in DB
	createPollParams := queries.CreatePollParams{
		Title: pgtype.Text{
			String: request.Title,
			Valid:  true,
		},
		DiscordID: pgtype.Text{
			String: request.DiscordId,
			Valid:  true,
		},
		DiscordAuthorID: pgtype.Text{
			String: request.DiscordAuthorId,
			Valid:  true,
		},
		DiscordGuildID: pgtype.Text{
			String: request.DiscordGuildId,
			Valid:  true,
		},
		ChannelID: pgtype.Text{
			String: request.ChannelId,
			Valid:  true,
		},
	}

	poll, err := qr.CreatePoll(ctx, createPollParams)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.CreatePollResponse{}, errors.Join(err, rollbackErr)
	}

	// Create new record for each option in polls_options table in DB
	grpcOptions := make([]*polls.Option, len(pollOptions))
	for i, option := range pollOptions {
		if utf8.RuneCountInString(option) > 50 {
			rollbackErr := tx.Rollback(ctx)
			return &polls.CreatePollResponse{}, errors.Join(rollbackErr, errors.New("длина варианта ответа не должна превышать 50 символов"))
		}

		createOptionParams := queries.CreateOptionParams{}
		createOptionParams.Title = pgtype.Text{
			String: option,
			Valid:  true,
		}
		createOptionParams.PollID = pgtype.Int4{
			Int32: poll.ID,
			Valid: true,
		}

		pollOption, createOptionErr := qr.CreateOption(ctx, createOptionParams)
		if createOptionErr != nil {
			rollbackErr := tx.Rollback(ctx)
			return &polls.CreatePollResponse{}, errors.Join(createOptionErr, rollbackErr)
		}

		grpcPollOption := polls.Option{
			Title:      pollOption.Title.String,
			CustomId:   fmt.Sprintf("poll_%d_choice_%d", poll.ID, pollOption.ID),
			TotalVotes: 0,
		}

		grpcOptions[i] = &grpcPollOption
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return &polls.CreatePollResponse{}, commitErr
	}

	return &polls.CreatePollResponse{
		PollId:  poll.ID,
		Options: grpcOptions,
	}, nil
}
