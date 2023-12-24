package poll

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
)

type VoteButtonHandler struct {
	poll_id   int32
	option_id int32

	Log      logger.Logger
	Session  *discordgo.Session
	Database repository.DB
}

func (v *VoteButtonHandler) Handle(i any) {
	vote, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	var err error
	defer func() {
		if err != nil {
			v.Log.Error("error registering vote", slog.String("error", err.Error()))
			err := interactionRespondError(
				"Произошла ошибка при регистрации голоса",
				err, v.Session, vote,
			)
			if err != nil {
				v.Log.Error(
					"error editing an interaction",
					slog.String("error", err.Error()),
				)
			}

			return
		}

		err = interactionRespondSuccess(
			"Голос зарегистрирован",
			v.Session, vote,
		)
		if err != nil {
			v.Log.Error(
				"error editing an interaction",
				slog.String("error", err.Error()),
			)
		}
	}()

	ctx := context.Background()

	err = interactionRespondLoading(
		"Голос регистрируется...",
		v.Session, vote,
	)
	if err != nil {
		v.Log.Error(
			"error responding to interaction",
			slog.String("error", err.Error()),
		)

		return
	}

	user, err := v.Database.Users().GetUserByDiscordID(ctx, vote.Member.User.ID)
	if err != nil {
		return
	}

	poll, err := v.Database.Polls().GetPollWithDetails(ctx, int(v.poll_id))
	if err != nil {
		return
	}

	newVote, err := v.Database.Polls().TryAddVote(ctx, user.ID, poll.ID, int(v.option_id))
	if err != nil {
		return
	}

	poll.Votes = append(poll.Votes, *newVote)

	discordAuthor, err := v.Session.User(poll.Author.DiscordUserID)
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)

		return
	}

	for _, msg := range poll.Messages {
		pollEmbed := buildPollEmbed(
			poll,
			discordAuthor,
		)

		_, err = v.Session.ChannelMessageEditEmbeds(
			msg.DiscordChannelID,
			msg.DiscordMessageID,
			pollEmbed,
		)
		if err != nil {
			err = errors.Join(domain.ErrInternal, err)
			v.Log.Error("error editing poll message", slog.String("error", err.Error()))
		}
	}
}
