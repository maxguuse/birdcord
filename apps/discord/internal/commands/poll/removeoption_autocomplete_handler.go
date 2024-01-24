package poll

import "github.com/bwmarrin/discordgo"

func (h *Handler) removeOptionAutocomplete(
	i *discordgo.Interaction,
) {
	data := i.ApplicationCommandData()

	switch {
	case data.Options[0].Focused:
		h.Log.Debug("poll option autocomplete")
	case data.Options[1].Focused:
		h.Log.Debug("option option autocomplete")
	}
}
