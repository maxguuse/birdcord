package poll

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/disroute"
	"github.com/samber/lo"
)

func (h *Handler) removeOptionAutocomplete(ctx *disroute.Ctx) []*discordgo.ApplicationCommandOptionChoice {
	data := ctx.Interaction().ApplicationCommandData()

	focusedOption, ok := lo.Find(data.Options[0].Options, func(o *discordgo.ApplicationCommandInteractionDataOption) bool {
		return o.Focused
	})
	if !ok {
		h.logger.Error("Failed to find focused option")

		return nil
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	switch focusedOption.Name {
	case "poll":
		choices = h.autocompletePollList(ctx)
	case "option":
		choices = h.autocompleteOptionList(ctx)
	}

	return choices
}
