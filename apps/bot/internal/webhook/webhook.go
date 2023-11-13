package webhook

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func GetBirdWebhook(s *discordgo.Session, channelID string) *discordgo.Webhook {
	channelWebhooks, err := s.ChannelWebhooks(channelID)
	webhook, ok := lo.Find(channelWebhooks, func(webhook *discordgo.Webhook) bool {
		return webhook.Name == "Bird"
	})

	if !ok {
		webhook, err = s.WebhookCreate(channelID, "Bird", "")
		if err != nil {
			fmt.Println("Error creating webhook:", err) //TODO Replace with logger
			return nil
		}

		return webhook
	}

	return webhook
}
