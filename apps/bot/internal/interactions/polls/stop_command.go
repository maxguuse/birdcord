package polls

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/embeds"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/samber/lo"
	"strings"
)

var stop = &discordgo.ApplicationCommandOption{
	Name:        "stop",
	Description: "Остановить опрос",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "poll",
			Description:  "Опрос",
			Type:         discordgo.ApplicationCommandOptionInteger,
			Required:     true,
			Autocomplete: true,
		},
	},
}

func (p *Polls) handleStopCommand(
	s *discordgo.Session,
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	pollId := options["poll"].IntValue()

	response, err := p.client.StopPoll(context.Background(), &polls.StopPollRequest{
		PollId: int32(pollId),
	})
	if err != nil {
		fmt.Println("Error from gRPC StopPoll:", err) //TODO Replace with logger
		respondWithError(s, i, errors.Join(err, errors.New("внутренняя ошибка бота, сообщите гусю")))
		return
	}

	message, err := s.ChannelMessage(response.ChannelId, response.DiscordId)
	if err != nil {
		fmt.Println("Error getting poll message:", err) //TODO Replace with logger
		respondWithError(s, i, errors.New("сообщение с опросом не найдено"))
		return
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: response.ChannelId,
		ID:      response.DiscordId,
		Embed: embeds.PollResults(
			response.Title,
			message.Embeds[0].Description,
			message.Embeds[0].Footer.Text,
			message.Embeds[0].Author.Name,
			message.Embeds[0].Author.IconURL,
			strings.Join(lo.Map(response.Winners, func(w *polls.Option, _ int) string {
				return w.Title
			}), ", "),
			response.TotalVotes,
		),
		Components: make([]discordgo.MessageComponent, 0),
	})
	if err != nil {
		fmt.Println("Error editing poll:", err) //TODO Replace with logger
		p.processPollFailure(s, i, message, err)
		return
	}

	// Send interaction response if poll message sent successfully
	p.processPollInteractionResponse(s, i, message, int32(pollId), "Опрос остановлен.")
}
