package poll

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func (h *Handler) removeOptionAutocomplete(i *discordgo.Interaction, options optionsMap) (string, error) {
	data := i.ApplicationCommandData()

	focusedOption, ok := lo.Find(data.Options[0].Options, func(o *discordgo.ApplicationCommandInteractionDataOption) bool {
		return o.Focused
	})
	if !ok {
		return "", errors.New("there's no focused option")
	}

	var err error
	switch focusedOption.Name {
	case "poll":
		_, err = h.autocompletePollList(i, options)
	case "option":
		_, err = h.autocompleteOptionList(i, options)
	}

	return "", err
}
