package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/samber/lo"
)

type Service struct {
	db repository.DB
}

func New(db repository.DB) *Service {
	return &Service{
		db: db,
	}
}

type CreateRequest struct {
	GuildID string
	UserID  string
	Poll    Poll
}

type Poll struct {
	Title   string
	Options string
}

func (s *Service) Create(ctx context.Context, r *CreateRequest) (*domain.PollWithDetails, error) {
	optionsList, err := processPollOptions(r.Poll.Options)
	if err != nil {
		return nil, err
	}

	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	user, err := s.db.Users().GetUserByDiscordID(ctx, r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll, err := s.db.Polls().CreatePoll(
		ctx,
		r.Poll.Title,
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return poll, nil
}

type StopRequest struct {
	GuildID string
	UserID  string
	PollID  int64
}

type StopResponse struct {
	Poll    *domain.PollWithDetails
	Winners []string
}

func (s *Service) Stop(ctx context.Context, r *StopRequest) (*StopResponse, error) {
	optionsWithVotes := make(map[domain.PollOption]int)

	var repoErr *repository.NotFoundError
	poll, err := s.db.Polls().GetPollWithDetails(ctx, int(r.PollID)) // Pass Guild ID as well
	if errors.As(err, &repoErr) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if poll.Author.DiscordUserID != r.UserID {
		return nil, ErrNotAuthor
	}

	if poll.Guild.DiscordGuildID != r.GuildID {
		return nil, ErrNotFound
	}

	var maxVotes int = 0
	for _, option := range poll.Options {
		optionVotes := lo.CountBy(poll.Votes, func(v domain.PollVote) bool {
			return v.OptionID == option.ID
		})

		optionsWithVotes[option] = optionVotes

		if optionVotes > maxVotes {
			maxVotes = optionVotes
		}
	}

	winners := lo.FilterMap(poll.Options, func(o domain.PollOption, _ int) (string, bool) {
		return o.Title, optionsWithVotes[o] == maxVotes
	})

	err = s.db.Polls().UpdatePollStatus(ctx, int(r.PollID), false)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return &StopResponse{
		Poll:    poll,
		Winners: winners,
	}, nil
}

func processPollOptions(rawOptions string) ([]string, error) {
	optionsList := strings.Split(rawOptions, "|")
	if len(optionsList) < 2 || len(optionsList) > 25 {
		return nil, &domain.UsersideError{
			Msg: "Количество вариантов опроса должно быть от 2 до 25 включительно.",
		}
	}
	if lo.SomeBy(optionsList, func(o string) bool {
		return utf8.RuneCountInString(o) > 50 || utf8.RuneCountInString(o) < 1
	}) {
		return nil, &domain.UsersideError{
			Msg: "Длина варианта опроса не может быть больше 50 или меньше 1 символа.",
		}
	}

	return optionsList, nil
}
