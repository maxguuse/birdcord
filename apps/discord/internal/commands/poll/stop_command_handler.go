package poll

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (h *Handler) stopPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		err = helpers.InteractionResponseProcess(h.Session, i, "Опрос остановлен.", err)
		if err != nil {
			h.Log.Error("error editing an interaction response", slog.String("error", err.Error()))
		}
	}()

	ctx := context.Background()

	optionsWithVotes := make(map[domain.PollOption]int)

	pollId := options["poll"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		err = errors.Join(domain.ErrUserSide, domain.ErrNotAuthor)

		return
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		err = errors.Join(domain.ErrUserSide, domain.ErrWrongGuild)

		return
	}

	var maxVotes int = 0
	for _, option := range poll.Options {
		optionVotes := lo.Filter(poll.Votes, func(v domain.PollVote, _ int) bool {
			return v.OptionID == option.ID
		})
		optionVotesAmount := len(optionVotes)

		optionsWithVotes[option] = optionVotesAmount

		if optionVotesAmount > maxVotes {
			maxVotes = optionVotesAmount
		}
	}

	winners := make([]domain.PollOption, 0, len(poll.Options))

	for _, option := range poll.Options {
		if optionsWithVotes[option] == maxVotes {
			winners = append(winners, option)
		}
	}

	discordAuthor, err := h.Session.User(poll.Author.DiscordUserID)
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	winnersList := lo.Map(winners, func(option domain.PollOption, _ int) string {
		return option.Title
	})

	pollEmbed := buildPollEmbed(poll, discordAuthor)
	pollEmbed[0].Fields = append(pollEmbed[0].Fields, &discordgo.MessageEmbedField{
		Name:   "Победители",
		Value:  strings.Join(winnersList, ","),
		Inline: true,
	})

	for _, msg := range poll.Messages {
		_, err = h.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         msg.DiscordMessageID,
			Channel:    msg.DiscordChannelID,
			Embeds:     pollEmbed,
			Components: make([]discordgo.MessageComponent, 0),
		})
		if err != nil {
			err = errors.Join(domain.ErrInternal, err)
			h.Log.Error("error editing poll message", slog.String("error", err.Error()))
		}
	}

	err = h.Database.Polls().UpdatePollStatus(ctx, int(pollId), false)
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}
}
