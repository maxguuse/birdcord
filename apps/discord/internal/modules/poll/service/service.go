package service

import (
	"context"
	"errors"

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

func (s *Service) GetPoll(ctx context.Context, r *GetPollRequest) (*domain.PollWithDetails, error) {
	var repoErr *repository.NotFoundError
	poll, err := s.db.Polls().GetPollWithDetails(ctx, int(r.PollID))
	if errors.As(err, &repoErr) {
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
	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, r.GuildID)
	if err != nil {
		return nil, err
	}

	user, err := s.db.Users().GetUserByDiscordID(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	polls, err := s.db.Polls().GetActivePolls(ctx, guild.ID, user.ID)
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

	if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
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

	err = s.db.Polls().UpdatePollStatus(ctx, int(r.PollID), false)
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

	var repoErr *repository.NotFoundError
	poll, err := s.db.Polls().GetPollWithDetails(ctx, int(pollId))
	if errors.As(err, &repoErr) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
		return nil, err
	}

	if len(poll.Options) == 25 {
		return nil, ErrTooManyOptions
	}

	newOption, err := s.db.Polls().AddPollOption(ctx, int(pollId), r.Option)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll.Options = append(poll.Options, *newOption)

	return poll, nil
}

func (s *Service) RemoveOption(ctx context.Context, r *RemoveOptionRequest) (*domain.PollWithDetails, error) {
	pollId := r.PollID
	optionId := r.OptionID

	var repoErr *repository.NotFoundError
	poll, err := s.db.Polls().GetPollWithDetails(ctx, int(pollId))
	if errors.As(err, &repoErr) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if err := validatePollAuthor(poll, r.UserID, r.GuildID); err != nil {
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

	err = s.db.Polls().RemovePollOption(ctx, int(optionId))
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

	user, err := s.db.Users().GetUserByDiscordID(ctx, r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll, err := s.db.Polls().GetPollWithDetails(ctx, vote.PollId)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	newVote, err := s.db.Polls().TryAddVote(ctx, user.ID, poll.ID, vote.OptionId)
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
	_, err := s.db.Polls().CreatePollMessage(
		ctx, r.Message.ID, r.Message.ChannelID, r.PollID,
	)

	return err
}
