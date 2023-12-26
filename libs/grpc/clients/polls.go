package clients

import (
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPolls() polls.PollsClient {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return polls.NewPollsClient(conn)
}
