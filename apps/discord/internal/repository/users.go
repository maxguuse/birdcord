package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

type UsersRepository interface {
	GetUserByDiscordID(
		ctx context.Context,
		id string,
	) (*domain.User, error)
}

type usersRepository struct {
	q *postgres.DB
}

func NewUsersRepository(q *postgres.DB) UsersRepository {
	return &usersRepository{
		q: q,
	}
}

func (u *usersRepository) GetUserByDiscordID(
	ctx context.Context,
	id string,
) (*domain.User, error) {
	result := &domain.User{}

	err := u.q.Transaction(ctx, func(q *queries.Queries) error {
		user, err := q.GetUserByDiscordID(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = q.CreateUser(ctx, id)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return nil
		}

		result.ID = int(user.ID)
		result.DiscordUserID = user.DiscordUserID

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
