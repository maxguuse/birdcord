package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"os"
)

func (p *Polls) handleVoteButton(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	v *Vote,
) {
	voteResponse, err := p.pollsClient.Vote(context.Background(), &polls.VoteRequest{
		DiscordUserId: v.UserID,
		PollId:        v.PollID,
		OptionId:      v.OptionID,
	})
	if err != nil {
		fmt.Println("Error voting:", err)
		return
	}

	pollMessageData := buildPollMessageData(
		voteResponse.Title,
		fmt.Sprintf("%d", v.PollID),
		voteResponse.TotalVotes,
		voteResponse.Options,
		i.Member,
	)

	_, err = s.InteractionResponseEdit(&discordgo.Interaction{
		AppID: os.Getenv("APP_ID"),
		Token: voteResponse.DiscordToken,
	}, &discordgo.WebhookEdit{
		Embeds:     &pollMessageData.Embeds,
		Components: &pollMessageData.Components,
	})
	if err != nil {
		fmt.Println("Error editing poll:", err)
		return
	}
}
