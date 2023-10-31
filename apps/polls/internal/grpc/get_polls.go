package grpc

import (
	"context"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
)

func (p PollsServer) getPolls(ctx context.Context) (*polls.GetPollsResponse, error) {
	activePolls, err := p.qr.GetPolls(ctx)
	if err != nil {
		return nil, err
	}

	grpcPolls := make([]*polls.Poll, 0, len(activePolls))
	for _, poll := range activePolls {
		grpcPoll := polls.Poll{
			Title: poll.Title.String,
			Id:    poll.ID,
		}
		grpcPolls = append(grpcPolls, &grpcPoll)
	}

	/*TODO Add "active" column to polls table and send only those polls where "active" is true
	 * 		Set "active" column to false in stopPolls method */

	/*TODO Add "guild_id" column to polls table and send only those polls where "guild_id" is equal to requested guild
	 */

	return &polls.GetPollsResponse{
		Polls: grpcPolls,
	}, nil
}
