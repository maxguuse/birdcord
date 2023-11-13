package grpc

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p PollsServer) invalidatePoll(ctx context.Context, request *polls.InvalidatePollRequest) (*emptypb.Empty, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	qr := queries.New(tx)

	err = qr.RemoveVotedUsers(ctx, pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	})
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &emptypb.Empty{}, errors.Join(err, rollbackErr)
	}

	err = qr.RemovePollOptions(ctx, pgtype.Int4{
		Int32: request.PollId,
		Valid: true,
	})
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &emptypb.Empty{}, errors.Join(err, rollbackErr)
	}

	err = qr.RemovePoll(ctx, request.PollId)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return &emptypb.Empty{}, errors.Join(err, rollbackErr)
	}

	return &emptypb.Empty{}, tx.Commit(ctx)
}
