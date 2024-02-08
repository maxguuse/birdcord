package poll

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) addPollOption(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) error {
	ctx := context.Background()

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

	newOption, err := h.Database.Polls().AddPollOption(ctx, int(pollId), options["option"].StringValue())
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	poll.Options = append(poll.Options, *newOption)

	pollEmbed := buildPollEmbed(poll, i.Member.User)
	actionRows := h.buildActionRows(poll, i.ID)

	var wg *errgroup.Group
	for _, msg := range poll.Messages {
		msg := msg
		wg.Go(func() error {
			_, err := h.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         msg.DiscordMessageID,
				Channel:    msg.DiscordChannelID,
				Content:    new(string),
				Components: actionRows,
				Embeds:     pollEmbed,
			})

			return err
		})

	}

	if err = wg.Wait(); err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
