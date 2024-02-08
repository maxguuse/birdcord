package poll

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) startPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) (string, error) {
	ctx := context.Background()

	optionsList, err := processPollOptions(options["options"].StringValue())
	if err != nil {
		return "", err
	}

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return "", err
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return "", err
	}

	poll, err := h.Database.Polls().CreatePoll(
		ctx,
		options["title"].StringValue(),
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return "", err
	}

	err = h.sendPollMessage(ctx, i, poll, optionsList)
	if err != nil {
		return "", err
	}

	return "Опрос успешно создан.", nil
}
