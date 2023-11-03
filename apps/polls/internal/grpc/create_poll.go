package grpc

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func (p PollsServer) createPoll(ctx context.Context, request *polls.CreatePollRequest) (*polls.CreatePollResponse, error) {
	// Create new record in polls table in DB
	createPollParams := queries.CreatePollParams{}
	createPollParams.Title = pgtype.Text{
		String: request.Title,
		Valid:  true,
	}
	createPollParams.DiscordToken = pgtype.Text{
		String: request.DiscordToken,
		Valid:  true,
	}
	createPollParams.DiscordAuthorID = pgtype.Text{
		String: request.DiscordAuthorId,
		Valid:  true,
	}
	createPollParams.DiscordGuildID = pgtype.Text{
		String: request.DiscordGuildId,
		Valid:  true,
	}

	poll, err := p.qr.CreatePoll(ctx, createPollParams)
	if err != nil {
		return nil, err
	}

	// Create new record for each option in polls_options table in DB
	grpcOptions := make([]*polls.Option, 0, len(request.Options))
	for _, option := range request.Options {
		createOptionParams := queries.CreateOptionParams{}
		createOptionParams.Title = pgtype.Text{
			String: option,
			Valid:  true,
		}
		createOptionParams.PollID = pgtype.Int4{
			Int32: poll.ID,
			Valid: true,
		}

		pollOption, err := p.qr.CreateOption(ctx, createOptionParams)
		if err != nil {
			return nil, err
		}

		grpcPollOption := polls.Option{
			Title:      pollOption.Title.String,
			Id:         pollOption.ID,
			TotalVotes: 0,
		}

		grpcOptions = append(grpcOptions, &grpcPollOption)
	}

	return &polls.CreatePollResponse{
		PollId:  poll.ID,
		Options: grpcOptions,
	}, nil
}
