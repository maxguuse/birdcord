package grpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type PollsServer struct {
	polls.UnimplementedPollsServer

	qr *queries.Queries
}

func StartPollsServer(
	lc fx.Lifecycle,
	qr *queries.Queries,
) error {
	server := &PollsServer{
		qr: qr,
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

func (p PollsServer) GetPolls(ctx context.Context, _ *empty.Empty) (*polls.GetPollsResponse, error) {
	return p.getPolls(ctx)
}
