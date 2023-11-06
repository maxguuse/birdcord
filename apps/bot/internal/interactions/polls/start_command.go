package polls

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/embeds"
	"github.com/maxguuse/birdcord/libs/grpc/generated/polls"
	"strings"
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
	pollOptions := strings.Split(options["options"].StringValue(), "|")
	if len(pollOptions) > 25 {
		processInvalidOptionsCount(s, i)
		return
	}

	// Send message with poll Embed
	message, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "Опрос формируется...",
	})
	if err != nil {
		fmt.Println("Error sending poll message:", err) //TODO Replace with logger
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
		return
	}

	// Create record in DB
	createPollResponse, err := p.client.CreatePoll(
		context.Background(),
		&polls.CreatePollRequest{
			Title:           options["title"].StringValue(),
			Options:         pollOptions,
			DiscordId:       message.ID,
			ChannelId:       message.ChannelID,
			DiscordAuthorId: i.Member.User.ID,
			DiscordGuildId:  i.GuildID,
		},
	)
	if err != nil {
		fmt.Println("Error creating poll:", err) //TODO Replace with logger
		processFailPoll(s, i, message)
		return
	}

	description, buttons := generateButtonsAndDescription(pollOptions, createPollResponse)

	actionRows := generateActionRows(buttons)

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: message.ChannelID,
		ID:      message.ID,
		Content: new(string),
		Embed: embeds.ActivePoll(
			options["title"].StringValue(),
			description,
			fmt.Sprintf("Poll ID: %d", createPollResponse.PollId),
			0,
			i.Member.Nick,
			i.Member.AvatarURL("1024"),
		),
		Components: actionRows,
	})
	if err != nil {
		fmt.Println("Error editing poll message:", err) //TODO Replace with logger
		processFailPoll(s, i, message)
		return
	}
}

func generateButtonsAndDescription(pollOptions []string, createPollResponse *polls.CreatePollResponse) (string, []discordgo.MessageComponent) {
	var description string
	buttons := make([]discordgo.MessageComponent, len(pollOptions))
	for i, option := range createPollResponse.Options {
		description += fmt.Sprintf("%d) %s > %d \n", i+1, option.Title, option.TotalVotes)
		buttons[i] = discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("poll_%d_choice_%d", createPollResponse.PollId, option.Id),
		}
	}
	return description, buttons
}

func generateActionRows(buttons []discordgo.MessageComponent) []discordgo.MessageComponent {
	actionRows := make([]discordgo.MessageComponent, 0, (len(buttons)+4)/5)
	for i := 0; i < len(buttons); i += 5 {
		actionRow := discordgo.ActionsRow{}
		for j := 0; j < 5 && i+j < len(buttons); j++ {
			actionRow.Components = append(actionRow.Components, buttons[i+j])
		}
		actionRows = append(actionRows, &actionRow)
	}
	return actionRows
}

func processInvalidOptionsCount(s *discordgo.Session, i *discordgo.Interaction) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Слишком много опций.",
		},
	})
	if err != nil {
		fmt.Println("Error responding to poll:", err) //TODO Replace with logger
	}
}

func processFailPoll(s *discordgo.Session, i *discordgo.Interaction, message *discordgo.Message) {
	deleteErr := s.ChannelMessageDelete(message.ChannelID, message.ID)
	if deleteErr != nil {
		fmt.Println("Error deleting poll message:", deleteErr) //TODO Replace with logger
		return
	}
	_, followupErr := s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{})
	if followupErr != nil {
		fmt.Println("Error sending followup to poll:", followupErr) //TODO Replace with logger
		return
	}
}
