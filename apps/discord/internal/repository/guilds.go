package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

type GuildsRepository interface {
	GetGuildByDiscordID(
		ctx context.Context,
		id string,
	) (*domain.Guild, error)
}

type guildsRepository struct {
	q *postgres.DB
}

func NewGuildsRepository(q *postgres.DB) GuildsRepository {
	return &guildsRepository{
		q: q,
	}
}

func (g *guildsRepository) GetGuildByDiscordID(
	ctx context.Context,
	id string,
) (*domain.Guild, error) {
	result := &domain.Guild{}

	err := g.q.Transaction(ctx, func(q *queries.Queries) error {
		guild, err := q.GetGuildByDiscordID(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			guild, err = q.CreateGuild(ctx, id)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		result.ID = int(guild.ID)
		result.DiscordGuildID = guild.DiscordGuildID

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
