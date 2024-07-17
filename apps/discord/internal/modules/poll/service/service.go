package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/poll/repository"
	"github.com/samber/lo"
)

type Service struct {
	repo repository.Repository
}

func New(
	repo repository.Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetPoll(ctx context.Context, r *GetPollRequest) (*domain.PollWithDetails, error) {
	poll, err := s.repo.GetPollWithDetails(ctx, int(r.PollID))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, int64(userId), int64(guildId)); err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *Service) GetActivePolls(ctx context.Context, r *GetActivePollsRequest) ([]*domain.Poll, error) {
	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
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
	optionsWithVotes := make(map[domain.PollOption]int)

	poll, err := s.repo.GetPollWithDetails(ctx, int(r.PollID)) // Pass Guild ID as well
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, int64(userId), int64(guildId)); err != nil {
		return nil, err
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

	err = s.repo.UpdatePollStatus(ctx, int(r.PollID), false)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return &StopResponse{
		Poll:    poll,
		Winners: winners,
	}, nil
}

func (s *Service) AddOption(ctx context.Context, r *AddOptionRequest) (*domain.PollWithDetails, error) {
	pollId := r.PollID

	poll, err := s.repo.GetPollWithDetails(ctx, int(pollId))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, int64(userId), int64(guildId)); err != nil {
		return nil, err
	}

	if len(poll.Options) == 25 {
		return nil, ErrTooManyOptions
	}

	newOption, err := s.repo.AddPollOption(ctx, int(pollId), r.Option)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll.Options = append(poll.Options, *newOption)

	return poll, nil
}

func (s *Service) RemoveOption(ctx context.Context, r *RemoveOptionRequest) (*domain.PollWithDetails, error) {
	pollId := r.PollID
	optionId := r.OptionID

	poll, err := s.repo.GetPollWithDetails(ctx, int(pollId))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, int64(userId), int64(guildId)); err != nil {
		return nil, err
	}

	optionVotes := lo.CountBy(poll.Votes, func(v domain.PollVote) bool {
		return v.OptionID == int(optionId)
	})

	if optionVotes > 0 {
		return nil, ErrOptionHasVotes
	}

	if len(poll.Options) <= 2 {
		return nil, ErrTooFewOptions
	}

	err = s.repo.RemovePollOption(ctx, int(optionId))
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll.Options = lo.Filter(poll.Options, func(o domain.PollOption, _ int) bool {
		return o.ID != int(optionId)
	})

	return poll, nil
}

func (s *Service) AddVote(ctx context.Context, r *AddVoteRequest) (*domain.PollWithDetails, error) {
	vote, err := parseVoteData(r.CustomID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll, err := s.repo.GetPollWithDetails(ctx, vote.PollId)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	userId, err := strconv.Atoi(r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	newVote, err := s.repo.TryAddVote(ctx, int64(userId), poll.ID, vote.OptionId)
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, &domain.UsersideError{
			Msg: "Вы уже проголосовали в этом опросе.",
		}
	}

	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll.Votes = append(poll.Votes, *newVote)

	return poll, nil
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
