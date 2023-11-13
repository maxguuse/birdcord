package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func (p PollsServer) stopPoll(ctx context.Context, request *polls.StopPollRequest) (*polls.StopPollResponse, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	qr := queries.New(tx)

	options, err := qr.GetOptionsWithVotesCount(ctx, pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	})
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.StopPollResponse{}, errors.Join(err, rollbackErr)
	}

	// Calculate winners
	var maxVotes int32 = 0
	for _, option := range options {
		if option.VoteCount > maxVotes {
			maxVotes = option.VoteCount
		}
	}

	winners := make([]*polls.Option, 0)
	for _, option := range options {
		if option.VoteCount == maxVotes {
			winners = append(winners, &polls.Option{
				CustomId:   fmt.Sprintf("poll_%d_choice_%d", request.PollId, option.ID),
				Title:      option.Title.String,
				TotalVotes: option.VoteCount,
			})
		}
	}

	// Calculate total votes
	var totalVotes int32 = 0
	for _, option := range options {
		totalVotes += option.VoteCount
	}

	poll, err := qr.GetPoll(ctx, request.PollId)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.StopPollResponse{}, errors.Join(err, rollbackErr)
	}

	err = qr.StopPoll(ctx, request.PollId)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.StopPollResponse{}, errors.Join(err, rollbackErr)
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return &polls.StopPollResponse{}, commitErr
	}

	return &polls.StopPollResponse{
		DiscordId:  poll.DiscordID.String,
		ChannelId:  poll.ChannelID.String,
		Winners:    winners,
		TotalVotes: totalVotes,
		Title:      poll.Title.String,
	}, nil
}
