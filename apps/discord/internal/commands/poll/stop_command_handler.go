package poll

import (
	"context"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) stopPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) error {
	ctx := context.Background()

	optionsWithVotes := make(map[domain.PollOption]int)

	pollId := options["poll"].IntValue()

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(pollId))
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != i.Member.User.ID {
		return errors.Join(domain.ErrUserSide, domain.ErrNotAuthor)
	}

	if poll.Guild.DiscordGuildID != i.GuildID {
		return errors.Join(domain.ErrUserSide, domain.ErrWrongGuild)
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
		return errors.Join(domain.ErrInternal, err)
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

	var wg *errgroup.Group
	for _, msg := range poll.Messages {
		msg := msg
		wg.Go(func() error {
			_, err = h.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         msg.DiscordMessageID,
				Channel:    msg.DiscordChannelID,
				Embeds:     pollEmbed,
				Components: make([]discordgo.MessageComponent, 0),
			})

			return err
		})
	}
	if err = wg.Wait(); err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = h.Database.Polls().UpdatePollStatus(ctx, int(pollId), false)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
