package polls

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func parseVoteFromButtonInteraction(i *discordgo.InteractionCreate) *Vote {
	customId := i.MessageComponentData().CustomID
	customIdParts := strings.Split(customId, "_")

	if len(customIdParts) != 4 {
		fmt.Println("Error parsing CustomID: ", customId, "len(customIdParts) != 4, invalid format")
		return nil
	}

	pollId, err := strconv.Atoi(customIdParts[1])
	if err != nil {
		fmt.Println("Error parsing CustomID: ", customId, err)
		return nil
	}
	optionId, err := strconv.Atoi(customIdParts[3])
	if err != nil {
		fmt.Println("Error parsing CustomID: ", customId, err)
		return nil
	}

	return &Vote{
		PollID:   int32(pollId),
		OptionID: int32(optionId),
		UserID:   i.Member.User.ID,
	}
}
