package tempvoice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/embed"
	"github.com/maxguuse/disroute"
)

func (h *Handler) setup(ctx *disroute.Ctx) disroute.Response {
	/* TODO

	Get []domain.TempvoiceHub from DB
	Get all channels from guild
	Filter all channels by domain.TempvoiceHub.DiscordChannelId to leave only hubs

	*/

	emb, comp := embed.HubsManagement(make([]*discordgo.Channel, 0))

	return disroute.Response{
		CustomResponse: &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
			Data: &discordgo.InteractionResponseData{
				Embeds:     emb,
				Components: comp,
			},
		},
	}
}
