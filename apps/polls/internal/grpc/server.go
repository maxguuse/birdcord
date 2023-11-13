package grpc

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type PollsServer struct {
	polls.UnimplementedPollsServer

	pool *pgxpool.Pool
}

func StartPollsServer(
	lc fx.Lifecycle,
	pool *pgxpool.Pool,
) error {
	server := &PollsServer{
		pool: pool,
	}

	grpcNetListener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	polls.RegisterPollsServer(grpcServer, server)
	reflection.Register(grpcServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := grpcServer.Serve(grpcNetListener); err != nil {
					panic(err)
				}
				fmt.Println("Server started")
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			fmt.Println("Server stopped")
			return nil
		},
	})

	return nil
}

func (p PollsServer) CreatePoll(ctx context.Context, request *polls.CreatePollRequest) (*polls.CreatePollResponse, error) {
	return p.createPoll(ctx, request)
}

func (p PollsServer) Vote(ctx context.Context, request *polls.VoteRequest) (*polls.VoteResponse, error) {
	return p.vote(ctx, request)
}

func (p PollsServer) GetActivePolls(ctx context.Context, request *polls.GetActivePollsRequest) (*polls.GetActivePollsResponse, error) {
	return p.getActivePolls(ctx, request)
}

func (p PollsServer) StopPoll(ctx context.Context, request *polls.StopPollRequest) (*polls.StopPollResponse, error) {
	return p.stopPoll(ctx, request)
}

func (p PollsServer) InvalidatePoll(ctx context.Context, request *polls.InvalidatePollRequest) (*emptypb.Empty, error) {
	return p.invalidatePoll(ctx, request)
}
