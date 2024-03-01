package poll

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
)

func (h *Handler) VoteBtnHandler(i *discordgo.Interaction) (string, error) {
	ctx := context.Background()

	vote, err := parseVoteData(i.MessageComponentData().CustomID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	user, err := h.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	poll, err := h.Database.Polls().GetPollWithDetails(ctx, vote.PollId)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	newVote, err := h.Database.Polls().TryAddVote(ctx, user.ID, poll.ID, vote.OptionId)
	if errors.Is(err, repository.ErrAlreadyExists) {
		return "", &domain.UsersideError{
			Msg: "Вы уже проголосовали в этом опросе.",
		}
	}

	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	poll.Votes = append(poll.Votes, *newVote)

	err = h.updatePollMessages(&UpdatePollMessageData{
		poll:        poll,
		interaction: i,
	})

	return "Голос зарегистрирован.", err
}

type voteData struct {
	PollId   int
	OptionId int
}

func parseVoteData(customId string) (*voteData, error) {
	parts := strings.Split(customId, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid custom id")
	}

	blob := parts[1]

	parts = strings.Split(blob, "_")
	if len(parts) != 6 {
		return nil, errors.New("invalid custom id")
	}

	poll_id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	option_id, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	return &voteData{
		PollId:   poll_id,
		OptionId: option_id,
	}, nil
}
