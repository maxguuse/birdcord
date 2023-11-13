package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

func (p PollsServer) getActivePolls(ctx context.Context, request *polls.GetActivePollsRequest) (*polls.GetActivePollsResponse, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return &polls.GetActivePollsResponse{}, err
	}
	qr := queries.New(tx)

	activePolls, err := qr.GetActivePolls(ctx, pgtype.Text{
		String: request.DiscordGuildId,
		Valid:  true,
	})
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &polls.GetActivePollsResponse{}, errors.Join(err, rollbackErr)
	}

	grpcPolls := make([]*polls.Poll, 0, len(activePolls))
	for _, poll := range activePolls {
		fmt.Println("Active poll:", poll.Title, poll.ID)
		grpcPoll := polls.Poll{
			Title: poll.Title.String,
			Id:    poll.ID,
		}
		grpcPolls = append(grpcPolls, &grpcPoll)
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return &polls.GetActivePollsResponse{}, commitErr
	}

	return &polls.GetActivePollsResponse{
		Polls: grpcPolls,
	}, nil
}
