package poll

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (h *Handler) removeOptionAutocomplete(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	data := i.ApplicationCommandData()

	focusedOption, ok := lo.Find(data.Options[0].Options, func(o *discordgo.ApplicationCommandInteractionDataOption) bool {
		return o.Focused
	})
	if !ok {
		return
	}

	h.Log.Debug("focused option", slog.Any("option", focusedOption))

	switch focusedOption.Name {
	case "poll":
		h.autocompletePollList(i, options)
	case "option":
		h.autocompleteOptionList(i, options)
	}
}