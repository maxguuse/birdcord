package grpc

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
)

func (p PollsServer) getActivePolls(ctx context.Context, request *polls.GetActivePollsRequest) (*polls.GetActivePollsResponse, error) {
	activePolls, err := p.qr.GetActivePolls(ctx, pgtype.Text{
		String: request.DiscordGuildId,
		Valid:  true,
	})
	if err != nil {
		return nil, err
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

	/*TODO Add "active" column to polls table and send only those polls where "active" is true
	 * 		Set "active" column to false in stopPolls method */

	/*TODO Add "guild_id" column to polls table and send only those polls where "guild_id" is equal to requested guild
	 */

	return &polls.GetActivePollsResponse{
		Polls: grpcPolls,
	}, nil
}
