package poll

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (p *CommandHandler) stopPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		if err != nil {
			p.Log.Error("error creating poll", slog.String("error", err.Error()))
			err := interactionRespondError(
				"Произошла ошибка при остановке опроса",
				err, p.Session, i,
			)
			if err != nil {
				p.Log.Error(
					"error editing an interaction",
					slog.String("error", err.Error()),
				)
			}

			return
		}

		err = interactionRespondSuccess(
			"Опрос остановлен!",
			p.Session, i,
		)
		if err != nil {
			p.Log.Error(
				"error editing an interaction",
				slog.String("error", err.Error()),
			)
		}
	}()

	_ = context.Background()

	err = interactionRespondLoading("Опрос останавливается...", p.Session, i)
	if err != nil {
		p.Log.Error(
			"error responding to interaction",
			slog.String("error", err.Error()),
		)
		return
	}

}
