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
	if !voteResponse.Success {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Вы уже участвовали в опросе \"%s\"", voteResponse.Title),
			},
		})
		if err != nil {
			fmt.Println("Error responding to vote interaction:", err)
		}
		return
	}

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

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Вы проголосовали в опросе \"%s\"", voteResponse.Title),
		},
	})
	if err != nil {
		fmt.Println("Error responding to vote interaction:", err)
		return
	}
}
