package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/embeds"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"github.com/samber/lo"
)

var start = &discordgo.ApplicationCommandOption{
	Name:        "start",
	Description: "Начать опрос",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "title",
			Description: "Заголовок опроса",
			Type:        discordgo.ApplicationCommandOptionString,
			MaxLength:   50,
			Required:    true,
		},
		{
			Name:        "options",
			Description: "Варианты ответа (разделите их символом '|')",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func (p *Polls) handleStart(
	s *discordgo.Session, i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	// Send message
	message, err := s.ChannelMessageSend(i.ChannelID, "Опрос формируется...")
	if err != nil {
		fmt.Println("Error sending poll message:", err) //TODO Replace with logger
		return
	}
	// Process poll data
	response, err := p.client.CreatePoll(
		context.Background(),
		&polls.CreatePollRequest{
			Title:           options["title"].StringValue(),
			Options:         options["options"].StringValue(),
			DiscordId:       message.ID,
			ChannelId:       message.ChannelID,
			DiscordAuthorId: i.Member.User.ID,
			DiscordGuildId:  i.GuildID,
		},
	)
	if err != nil {
		fmt.Println("Error from gRPC CreatePoll:", err) //TODO Replace with logger
		processPollFailure(s, i, message, err)
		return
	}
	// Edit poll message with actual data
	components := generatePollComponents(response)
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: message.ChannelID,
		ID:      message.ID,
		Content: new(string),
		Embed: embeds.ActivePoll(
			options["title"].StringValue(),
			components.description,
			fmt.Sprintf("Poll ID: %d", response.PollId),
			0,
			i.Member.Nick,
			i.Member.AvatarURL("1024"),
		),
		Components: components.actionRows,
	})
	if err != nil {
		fmt.Println("Error editing poll message:", err) //TODO Replace with logger
		processPollFailure(s, i, message, err)
		return
	}
	// Send interaction response if poll message sent successfully
	err = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Опрос успешно создан.",
		},
	})
	if err != nil {
		fmt.Println("Error responding to poll:", err) //TODO Replace with logger
		processPollFailure(s, i, message, err)
	}
}

type pollComponents struct {
	description string
	actionRows  []discordgo.MessageComponent
}

func generatePollComponents(response *polls.CreatePollResponse) *pollComponents {
	description := ""
	actionRows := make([]*discordgo.ActionsRow, 0, (len(response.Options)+4)/5)

	for i, option := range response.Options {
		description += fmt.Sprintf("%d) %s > %d \n", i+1, option.Title, option.TotalVotes)

		// If the current index is a multiple of 5, create a new ActionsRow
		if i%5 == 0 {
			actionRow := &discordgo.ActionsRow{
				Components: make([]discordgo.MessageComponent, 0, 5),
			}
			actionRows = append(actionRows, actionRow)
		}

		// Add the current button to the last ActionsRow
		lastActionRow := actionRows[len(actionRows)-1]
		lastActionRow.Components = append(lastActionRow.Components, discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: option.CustomId,
		})
	}

	return &pollComponents{
		description: description,
		actionRows: lo.Map(actionRows, func(actionRow *discordgo.ActionsRow, _ int) discordgo.MessageComponent {
			return actionRow
		}),
	}
}

func processPollFailure(s *discordgo.Session, i *discordgo.Interaction, message *discordgo.Message, err error) {
	deleteErr := s.ChannelMessageDelete(message.ChannelID, message.ID)
	if deleteErr != nil {
		fmt.Println("Error deleting poll message:", deleteErr) //TODO Replace with logger
		return
	}

	responseErr := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if responseErr != nil {
		fmt.Println("Error responding to poll interaction with error:", responseErr) //TODO Replace with logger
		return
	}

	_, followupErr := s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
		Content: "Ошибка при создании опроса!",
		Flags:   discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Сообщение об ошибке",
				Description: err.Error(),
			},
		},
	})
	if followupErr != nil {
		fmt.Println("Error sending poll error message:", followupErr) //TODO Replace with logger
		return
	}
}
