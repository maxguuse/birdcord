package client

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) componentKeyFunc(i *discordgo.Interaction) string {
	parts := strings.Split(i.MessageComponentData().CustomID, ":")

	return parts[0]
}
