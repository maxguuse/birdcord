package poll

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/commands/helpers"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
)

func (h *Handler) BuildVoteButtonHandler(poll_id, option_id int32) func(*discordgo.Interaction) {
	return func(i *discordgo.Interaction) {
		var err error
		defer func() {
			err = helpers.InteractionResponseProcess(h.Session, i, "Голос зарегистрирован.", err)
			if err != nil {
				h.Log.Error("error editing an interaction response", slog.String("error", err.Error()))
			}
		}()

		ctx := context.Background()

		user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
		if err != nil {
			err = errors.Join(domain.ErrInternal, err)

			return
		}

		poll, err := h.Database.Polls().GetPollWithDetails(ctx, int(poll_id))
		if err != nil {
			err = errors.Join(domain.ErrInternal, err)

			return
		}

		newVote, err := h.Database.Polls().TryAddVote(ctx, user.ID, poll.ID, int(option_id))
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				err = &domain.UsersideError{
					Msg: "Вы уже проголосовали в этом опросе.",
				}
			} else if err != nil {
				err = errors.Join(domain.ErrInternal, err)
			}

			return
		}

		poll.Votes = append(poll.Votes, *newVote)

		err = h.updatePollMessages(&UpdatePollMessageData{
			poll:        poll,
			interaction: i,
		})
	}
}
