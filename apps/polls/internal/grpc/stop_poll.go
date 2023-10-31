package grpc

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
)

func (p PollsServer) StopPoll(ctx context.Context, request *polls.StopPollRequest) (*polls.StopPollResponse, error) {
	options, err := p.qr.GetOptionsWithVotesCount(ctx, pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	})
	if err != nil {
		return nil, err
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
				Id:         option.ID,
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

	poll, err := p.qr.GetPoll(ctx, request.PollId)
	if err != nil {
		return nil, err
	}

	return &polls.StopPollResponse{
		DiscordToken: poll.DiscordToken.String,
		Winners:      winners,
		TotalVotes:   totalVotes,
		Title:        poll.Title.String,
	}, nil
}
