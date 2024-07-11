package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
)

type PollsRepository interface {
	TryAddVote(
		ctx context.Context,
		userID, pollID, optionID int,
	) (*domain.PollVote, error)
	CreatePollMessage(
		ctx context.Context,
		discordMessageId, discordChannelId string,
		pollId int,
	) (*domain.PollMessage, error)
	UpdatePollStatus(
		ctx context.Context,
		pollId int,
		isActive bool,
	) error
	GetActivePolls(
		ctx context.Context,
		guildID int,
		authorID int,
	) ([]*domain.Poll, error)
	AddPollOption(
		ctx context.Context,
		pollID int,
		pollOption string,
	) (*domain.PollOption, error)
	RemovePollOption(
		ctx context.Context,
		optionID int,
	) error
}

type pollsRepository struct {
	q *postgres.DB
}

func NewPollsRepository(q *postgres.DB) PollsRepository {
	return &pollsRepository{
		q: q,
	}
}

func (p *pollsRepository) GetActivePolls(
	ctx context.Context,
	guildID int,
	authorID int,
) ([]*domain.Poll, error) {
	result := []*domain.Poll{}

	err := p.q.Transaction(ctx, func(q *queries.Queries) error {
		polls, err := q.GetActivePolls(ctx, queries.GetActivePollsParams{
			GuildID: int32(guildID),
			AuthorID: pgtype.Int4{
				Int32: int32(authorID),
				Valid: true,
			},
		})
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPollNotFound
		} else if err != nil {
			return err
		}

		result = lo.Map(polls, func(p queries.Poll, _ int) *domain.Poll {
			return &domain.Poll{
				ID:        int(p.ID),
				Title:     p.Title,
				CreatedAt: p.CreatedAt.Time,
			}
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pollsRepository) TryAddVote(
	ctx context.Context,
	userID, pollID, optionID int,
) (*domain.PollVote, error) {
	result := &domain.PollVote{}

	err := p.q.Transaction(ctx, func(q *queries.Queries) error {
		vote, err := q.AddVote(ctx, queries.AddVoteParams{
			UserID:   int32(userID),
			PollID:   int32(pollID),
			OptionID: int32(optionID),
		})

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
		if err != nil {
			return err
		}

		result.ID = int(vote.ID)
		result.OptionID = int(vote.OptionID)
		result.UserID = int(vote.UserID)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pollsRepository) CreatePollMessage(
	ctx context.Context,
	discordMessageId, discordChannelId string,
	pollId int,
) (*domain.PollMessage, error) {
	result := &domain.PollMessage{}

	err := p.q.Transaction(ctx, func(q *queries.Queries) error {
		msg, err := q.CreateMessage(ctx, queries.CreateMessageParams{
			DiscordMessageID: discordMessageId,
			DiscordChannelID: discordChannelId,
		})
		if err != nil {
			return err
		}

		pollMsg, err := q.CreatePollMessage(ctx, queries.CreatePollMessageParams{
			PollID:    int32(pollId),
			MessageID: msg.ID,
		})
		if err != nil {
			return err
		}

		result.ID = int(pollMsg.ID)
		result.MessageID = int(msg.ID)
		result.DiscordMessageID = msg.DiscordMessageID
		result.DiscordChannelID = msg.DiscordChannelID

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pollsRepository) UpdatePollStatus(
	ctx context.Context,
	pollId int,
	status bool,
) error {
	err := p.q.Transaction(ctx, func(q *queries.Queries) error {
		err := q.UpdatePollStatus(ctx, queries.UpdatePollStatusParams{
			ID:       int32(pollId),
			IsActive: status,
		})

		return err
	})

	return err
}

func (p *pollsRepository) AddPollOption(
	ctx context.Context,
	pollId int,
	pollOption string,
) (*domain.PollOption, error) {
	result := &domain.PollOption{}

	err := p.q.Transaction(ctx, func(q *queries.Queries) error {
		newOption, err := q.CreatePollOption(ctx, queries.CreatePollOptionParams{
			Title:  pollOption,
			PollID: int32(pollId),
		})
		if err != nil {
			return err
		}

		result.ID = int(newOption.ID)
		result.Title = newOption.Title

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pollsRepository) RemovePollOption(
	ctx context.Context,
	optionId int,
) error {
	err := p.q.Transaction(ctx, func(q *queries.Queries) error {
		err := q.DeletePollOption(ctx, int32(optionId))

		return err
	})

	return err
}
