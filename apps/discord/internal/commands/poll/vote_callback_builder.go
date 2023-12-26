package poll

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"go.uber.org/fx"
)

type VoteCallbackBuilder struct {
	Log      logger.Logger
	Session  *discordgo.Session
	Database repository.DB
}

type VoteCallbackBuilderOpts struct {
	fx.In

	Log      logger.Logger
	Session  *discordgo.Session
	Database repository.DB
}

func NewVoteCallbackBuilder(opts VoteCallbackBuilderOpts) *VoteCallbackBuilder {
	return &VoteCallbackBuilder{
		Log:      opts.Log,
		Session:  opts.Session,
		Database: opts.Database,
	}
}

type Callback = func(i *discordgo.Interaction)

func (h *VoteCallbackBuilder) Build(poll_id, option_id int32) Callback {
	return func(i *discordgo.Interaction) {
		var err error
		defer func() {
			if err != nil {
				h.Log.Error("error registering vote", slog.String("error", err.Error()))
				err := interactionRespondError("Произошла ошибка при регистрации голоса", err, h.Session, i)
				if err != nil {
					h.Log.Error("error editing an interaction", slog.String("error", err.Error()))
				}

				return
			}

			err = interactionRespondSuccess("Голос зарегистрирован", h.Session, i)
			if err != nil {
				h.Log.Error("error editing an interaction", slog.String("error", err.Error()))
			}
		}()

		err = interactionRespondLoading("Голос регистрируется...", h.Session, i)
		if err != nil {
			h.Log.Error("error responding to interaction", slog.String("error", err.Error()))
		}

		ctx := context.Background()

		user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
		if err != nil {
			return
		}

		poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(poll_id))
		if err != nil {
			return
		}

		newVote, err := h.Database.Polls().TryAddVote(ctx, user.ID, poll.ID, int(option_id))
		if err != nil {
			return
		}

		poll.Votes = append(poll.Votes, *newVote)

		discordAuthor, err := h.Session.User(poll.Author.DiscordUserID)
		if err != nil {
			err = errors.Join(domain.ErrInternal, err)

			return
		}

		for _, msg := range poll.Messages {
			pollEmbed := buildPollEmbed(poll, discordAuthor)

			_, err = h.Session.ChannelMessageEditEmbeds(
				msg.DiscordChannelID,
				msg.DiscordMessageID,
				pollEmbed,
			)
			if err != nil {
				err = errors.Join(domain.ErrInternal, err)
				h.Log.Error("error editing poll message", slog.String("error", err.Error()))
			}
		}
	}
}