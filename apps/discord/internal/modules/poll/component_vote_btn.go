package poll

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/service"
	"github.com/maxguuse/disroute"
)

func (h *Handler) VoteBtnHandler(ctx *disroute.Ctx) disroute.Response {
	poll, err := h.service.AddVote(ctx.Context(), &service.AddVoteRequest{
		UserID:   ctx.Interaction().Member.User.ID,
		CustomID: ctx.Interaction().MessageComponentData().CustomID,
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll: poll,
		ctx:  ctx,
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	return disroute.Response{
		Message: "Голос добавлен.",
	}
}
