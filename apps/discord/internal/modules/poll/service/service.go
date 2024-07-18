package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/repository"
	"github.com/maxguuse/birdcord/libs/jet/pgerrors"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/samber/lo"
)

type Service struct {
	repo      repository.Repository
	txManager *txmanager.TxManager
}

func New(
	repo repository.Repository,
	txManager *txmanager.TxManager,
) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *Service) GetPoll(ctx context.Context, r *GetPollRequest) (*domain.PollWithDetails, error) {
	poll, err := s.repo.GetPollWithDetails(ctx, int(r.PollID))
	if errors.Is(err, pgerrors.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *Service) GetActivePolls(ctx context.Context, r *GetActivePollsRequest) ([]*domain.Poll, error) {
	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, err
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, err
	}

	polls, err := s.repo.GetActivePolls(ctx, int64(guildId), int64(userId))
	if err != nil {
		return nil, err
	}

	return polls, nil
}

func (s *Service) Create(ctx context.Context, r *CreateRequest) (*domain.PollWithDetails, error) {
	optionsList, err := processPollOptions(r.Poll.Options)
	if err != nil {
		return nil, err
	}

	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	authorId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll, err := s.repo.CreatePoll(
		ctx,
		int64(guildId),
		int64(authorId),
		r.Poll.Title,
		optionsList,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return poll, nil
}

func (s *Service) Stop(ctx context.Context, r *StopRequest) (*StopResponse, error) {
	resp := &StopResponse{}

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		poll, err := s.repo.GetPollWithDetails(ctx, int(r.PollID))
		if errors.Is(err, pgerrors.ErrNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
			return err
		}

		var maxVotes = 0
		var optionsWithVotes = make(map[domain.PollOption]int)
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

		err = s.repo.UpdatePollStatus(ctx, int(r.PollID), false)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		resp.Poll = poll
		resp.Winners = winners

		return nil
	})
	if errors.Is(err, txmanager.ErrTx) {
		return nil, errors.Join(domain.ErrInternal, err)
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) AddOption(ctx context.Context, r *AddOptionRequest) (*domain.PollWithDetails, error) {
	result := &domain.PollWithDetails{}

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		poll, err := s.repo.GetPollWithDetails(ctx, int(r.PollID))
		if errors.Is(err, pgerrors.ErrNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
			return err
		}

		if len(poll.Options) == 25 {
			return ErrTooManyOptions
		}

		newOption, err := s.repo.AddPollOption(ctx, int(r.PollID), r.Option)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		poll.Options = append(poll.Options, *newOption)

		result = poll

		return nil
	})
	if errors.Is(err, txmanager.ErrTx) {
		return nil, errors.Join(domain.ErrInternal, err)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) RemoveOption(ctx context.Context, r *RemoveOptionRequest) (*domain.PollWithDetails, error) {
	result := &domain.PollWithDetails{}

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		poll, err := s.repo.GetPollWithDetails(ctx, int(r.PollID))
		if errors.Is(err, pgerrors.ErrNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
			return err
		}

		optionVotes := lo.CountBy(poll.Votes, func(v domain.PollVote) bool {
			return v.OptionID == int(r.OptionID)
		})

		if optionVotes > 0 {
			return ErrOptionHasVotes
		}

		if len(poll.Options) <= 2 {
			return ErrTooFewOptions
		}

		err = s.repo.RemovePollOption(ctx, int(r.OptionID))
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		poll.Options = lo.Filter(poll.Options, func(o domain.PollOption, _ int) bool {
			return o.ID != int(r.OptionID)
		})

		result = poll

		return nil
	})
	if errors.Is(err, txmanager.ErrTx) {
		return nil, errors.Join(domain.ErrInternal, err)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) AddVote(ctx context.Context, r *AddVoteRequest) (*domain.PollWithDetails, error) {
	result := &domain.PollWithDetails{}

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		vote, err := parseVoteData(r.CustomID)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		poll, err := s.repo.GetPollWithDetails(ctx, vote.PollId)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		userId, err := strconv.Atoi(r.UserID)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		newVote, err := s.repo.TryAddVote(ctx, int64(userId), poll.ID, vote.OptionId)
		if errors.Is(err, pgerrors.ErrDuplicateKey) {
			return &domain.UsersideError{
				Msg: "Вы уже проголосовали в этом опросе.",
			}
		}

		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		poll.Votes = append(poll.Votes, *newVote)

		result = poll

		return nil
	})
	if errors.Is(err, txmanager.ErrTx) {
		return nil, errors.Join(domain.ErrInternal, err)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) CreateMessage(ctx context.Context, r *CreateMessageRequest) error {
	messageId, err := strconv.Atoi(r.Message.ID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	channelId, err := strconv.Atoi(r.Message.ChannelID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	_, err = s.repo.CreatePollMessage(
		ctx, int64(messageId), int64(channelId), r.PollID,
	)

	return err
}
