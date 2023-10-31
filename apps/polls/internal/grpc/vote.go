package grpc

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func (p PollsServer) vote(ctx context.Context, request *polls.VoteRequest) (*polls.VoteResponse, error) {
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

	_, err := p.qr.GetUserByIdForPoll(ctx, getUserByIdForPollParams)
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
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

	err = p.qr.AddVotedUser(ctx, addVotedUserParams)
	if err != nil {
		return nil, err
	}

	token, err := p.qr.GetToken(ctx, request.PollId)

	pollIdForPg := pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	}
	options, err := p.qr.GetOptionsWithVotesCount(ctx, pollIdForPg)
	if err != nil {
		return nil, err
	}

	grpcOptions := make([]*polls.Option, 0, len(options))
	var totalVotes int32 = 0
	for _, option := range options {
		grpcPollOption := polls.Option{
			Title:      option.Title.String,
			Id:         option.ID,
			TotalVotes: option.VoteCount,
		}
		grpcOptions = append(grpcOptions, &grpcPollOption)
		totalVotes += option.VoteCount
	}

	poll, err := p.qr.GetPoll(ctx, request.PollId)
	if err != nil {
		return nil, err
	}

	return &polls.VoteResponse{
		DiscordToken: token.String,
		Options:      grpcOptions,
		TotalVotes:   totalVotes,
		Title:        poll.Title.String,
	}, nil
}
