package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func (p PollsServer) vote(ctx context.Context, request *polls.VoteRequest) (*polls.VoteResponse, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return &polls.VoteResponse{}, err
	}
	qr := queries.New(tx)

	// Check if user has already voted for this poll
	getUserByIdForPollParams := queries.GetUserByIdForPollParams{}
	getUserByIdForPollParams.DiscordID = pgtype.Text{
		String: request.DiscordUserId,
		Valid:  true,
	}
	getUserByIdForPollParams.PollID = pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	}

	poll, err := qr.GetPoll(ctx, request.PollId)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.VoteResponse{}, errors.Join(err, rollbackErr)
	}

	_, err = qr.GetUserByIdForPoll(ctx, getUserByIdForPollParams)
	if !errors.Is(err, pgx.ErrNoRows) {
		return &polls.VoteResponse{
			Title:   poll.Title.String,
			Success: false,
		}, err
	}

	// If all conditions are met, add user to voted_users table
	addVotedUserParams := queries.AddVotedUserParams{}
	addVotedUserParams.DiscordID = pgtype.Text{
		String: request.DiscordUserId,
		Valid:  true,
	}
	addVotedUserParams.OptionID = pgtype.Int4{
		Int32: request.OptionId,
		Valid: true,
	}
	addVotedUserParams.PollID = pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	}

	err = qr.AddVotedUser(ctx, addVotedUserParams)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.VoteResponse{}, errors.Join(err, rollbackErr)
	}

	token, err := qr.GetToken(ctx, request.PollId)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.VoteResponse{}, errors.Join(err, rollbackErr)
	}

	pollIdForPg := pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	}
	options, err := qr.GetOptionsWithVotesCount(ctx, pollIdForPg)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.VoteResponse{}, errors.Join(err, rollbackErr)
	}

	grpcOptions := make([]*polls.Option, 0, len(options))
	var totalVotes int32 = 0
	for _, option := range options {
		grpcPollOption := polls.Option{
			Title:      option.Title.String,
			CustomId:   fmt.Sprintf("poll_%d_choice_%d", request.PollId, option.ID),
			TotalVotes: option.VoteCount,
		}
		grpcOptions = append(grpcOptions, &grpcPollOption)
		totalVotes += option.VoteCount
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return &polls.VoteResponse{}, commitErr
	}

	return &polls.VoteResponse{
		DiscordToken: token.String,
		Options:      grpcOptions,
		TotalVotes:   totalVotes,
		Title:        poll.Title.String,
		Success:      true,
	}, nil
}
