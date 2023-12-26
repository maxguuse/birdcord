package poll

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
)

func (h *Handler) startPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		err = helpers.InteractionResponseProcess(h.Session, i, "Опрос создан.", err)
		if err != nil {
			h.Log.Error("error editing an interaction response", slog.String("error", err.Error()))
		}
	}()

	ctx := context.Background()

	optionsList, err := processPollOptions(options["options"].StringValue())
	if err != nil {
		return
	}

	guild, err := h.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return
	}

	poll, err := h.Database.Polls().CreatePoll(
		ctx,
		options["title"].StringValue(),
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return
	}

	err = h.sendPollMessage(ctx, i, poll, optionsList)
}
