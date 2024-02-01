package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
)

type PollsRepository interface {
	CreatePoll(
		ctx context.Context,
		title string,
		guildID, authorID int,
		pollOptions []string,
	) (*domain.PollWithDetails, error)
	GetPollWithDetails(
		ctx context.Context,
		id int,
	) (*domain.PollWithDetails, error)
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

func (p *pollsRepository) CreatePoll(
	ctx context.Context,
	title string,
	guildID, authorID int,
	pollOptions []string,
) (*domain.PollWithDetails, error) {
	result := &domain.PollWithDetails{}

	err := p.q.Transaction(func(q *queries.Queries) error {
		poll, err := q.CreatePoll(ctx, queries.CreatePollParams{
			Title: title,
			AuthorID: pgtype.Int4{
				Int32: int32(authorID),
				Valid: true,
			},
			GuildID: int32(guildID),
		})
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.ID = int(poll.ID)
		result.Title = poll.Title
		result.CreatedAt = poll.CreatedAt.Time

		author, err := q.GetUserById(ctx, int32(authorID))
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Author = domain.PollAuthor{
			ID:            int(author.ID),
			DiscordUserID: author.DiscordUserID,
		}

		guild, err := q.GetGuildByID(ctx, int32(guildID))
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Guild = domain.PollGuild{
			ID:             int(guild.ID),
			DiscordGuildID: guild.DiscordGuildID,
		}

		options, err := q.CreatePollOptions(ctx, queries.CreatePollOptionsParams{
			Titles: pollOptions,
			PollID: int32(result.ID),
		})
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Options = lo.Map(options, func(r queries.CreatePollOptionsRow, _ int) domain.PollOption {
			return domain.PollOption{
				ID:    int(r.PollOption.ID),
				Title: r.PollOption.Title,
			}
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pollsRepository) GetPollWithDetails(
	ctx context.Context,
	id int,
) (*domain.PollWithDetails, error) {
	result := &domain.PollWithDetails{}

	err := p.q.Transaction(func(q *queries.Queries) error {
		poll, err := q.GetPoll(ctx, int32(id))
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.ID = int(poll.ID)
		result.Title = poll.Title
		result.CreatedAt = poll.CreatedAt.Time
		result.IsActive = poll.IsActive

		user, err := q.GetUserById(ctx, poll.AuthorID.Int32)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Author = domain.PollAuthor{
			ID:            int(user.ID),
			DiscordUserID: user.DiscordUserID,
		}

		guild, err := q.GetGuildByID(ctx, poll.GuildID)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Guild = domain.PollGuild{
			ID:             int(guild.ID),
			DiscordGuildID: guild.DiscordGuildID,
		}

		pollMessages, err := q.GetFullPollMessages(ctx, poll.ID)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Messages = lo.Map(
			pollMessages,
			func(m queries.GetFullPollMessagesRow, _ int) domain.PollMessage {
				return domain.PollMessage{
					ID:               int(m.ID),
					MessageID:        int(m.MessageID),
					DiscordMessageID: m.DiscordMessageID.String,
					DiscordChannelID: m.DiscordChannelID.String,
				}
			})

		pollOptions, err := q.GetPollOptions(ctx, poll.ID)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Options = lo.Map(pollOptions, func(o queries.PollOption, _ int) domain.PollOption {
			return domain.PollOption{
				ID:    int(o.ID),
				Title: o.Title,
			}
		})

		pollVotes, err := q.GetPollVotes(ctx, poll.ID)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		result.Votes = lo.Map(pollVotes, func(v queries.PollVote, _ int) domain.PollVote {
			return domain.PollVote{
				ID:       int(v.ID),
				OptionID: int(v.OptionID),
				UserID:   int(v.UserID),
			}
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pollsRepository) GetActivePolls(
	ctx context.Context,
	guildID int,
	authorID int,
) ([]*domain.Poll, error) {
	result := []*domain.Poll{}

	err := p.q.Transaction(func(q *queries.Queries) error {
		polls, err := q.GetActivePolls(ctx, queries.GetActivePollsParams{
			GuildID: int32(guildID),
			AuthorID: pgtype.Int4{
				Int32: int32(authorID),
				Valid: true,
			},
		})
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
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

	err := p.q.Transaction(func(q *queries.Queries) error {
		vote, err := q.AddVote(ctx, queries.AddVoteParams{
			UserID:   int32(userID),
			PollID:   int32(pollID),
			OptionID: int32(optionID),
		})

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.Join(
				domain.ErrUserSide,
				domain.ErrAlreadyVoted,
			)
		}
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
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

	err := p.q.Transaction(func(q *queries.Queries) error {
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
	err := p.q.Transaction(func(q *queries.Queries) error {
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

	err := p.q.Transaction(func(q *queries.Queries) error {
		newOption, err := q.CreatePollOption(ctx, queries.CreatePollOptionParams{
			Title:  pollOption,
			PollID: int32(pollId),
		})
		if err != nil {
			return errors.Join(
				domain.ErrInternal,
				err,
			)
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
	err := p.q.Transaction(func(q *queries.Queries) error {
		err := q.DeletePollOption(ctx, int32(optionId))

		return err
	})

	return err
}
