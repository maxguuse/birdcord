package tempvoice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/embed"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/service"
	"github.com/maxguuse/disroute"
)

func (h *Handler) setup(ctx *disroute.Ctx) disroute.Response {
	hubs := h.service.GetHubs(ctx.Context(), service.GetHubsRequest{
		DiscordGuildId: ctx.Interaction().GuildID,
		Session:        ctx.Session(),
		Interaction:    ctx.Interaction(),
	})

	emb, comp := embed.HubsManagement(hubs)

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
