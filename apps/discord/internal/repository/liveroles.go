package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
)

type LiverolesRepository interface {
	CreateLiverole(
		ctx context.Context,
		discordRoleId string,
		guildID int,
	) (*domain.Liverole, error)
	GetLiveroles(
		ctx context.Context,
		guildId int,
	) ([]*domain.Liverole, error)
	GetLiverole(
		ctx context.Context,
		guildID int,
		discordRoleID string,
	) (*domain.Liverole, error)
	DeleteLiverole(
		ctx context.Context,
		roleID int,
	) error
	DeleteLiveroles(
		ctx context.Context,
		guildID int,
		discordRolesIds []string,
	) error
}

type liverolesRepository struct {
	q *postgres.DB
}

func NewLiverolesRepository(q *postgres.DB) LiverolesRepository {
	return &liverolesRepository{
		q: q,
	}
}

func (l *liverolesRepository) CreateLiverole(
	ctx context.Context,
	discordRoleId string,
	guildID int,
) (*domain.Liverole, error) {
	result := &domain.Liverole{}

	err := l.q.Transaction(ctx, func(q *queries.Queries) error {
		role, err := q.CreateRole(ctx, queries.CreateRoleParams{
			GuildID:       int32(guildID),
			DiscordRoleID: discordRoleId,
		})
		if err != nil {
			return err
		}

		liverole, err := q.CreateLiveRole(ctx, role.ID)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
		if err != nil {
			return err
		}

		result.ID = int(liverole.ID)
		result.DiscordRoleID = role.DiscordRoleID
		result.GuildID = int(role.GuildID)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, err
}

func (l *liverolesRepository) DeleteLiverole(
	ctx context.Context,
	roleID int,
) error {
	err := l.q.Transaction(ctx, func(q *queries.Queries) error {
		err := q.DeleteRoleByID(ctx, int32(roleID))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (l *liverolesRepository) DeleteLiveroles(
	ctx context.Context,
	guildID int,
	discordRolesIds []string,
) error {
	err := l.q.Transaction(ctx, func(q *queries.Queries) error {
		err := q.DeleteRolesByGuildID(ctx, queries.DeleteRolesByGuildIDParams{
			GuildID:        int32(guildID),
			DiscordRoleIds: discordRolesIds,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (l *liverolesRepository) GetLiveroles(
	ctx context.Context,
	guildId int,
) ([]*domain.Liverole, error) {
	result := make([]*domain.Liverole, 0)

	err := l.q.Transaction(ctx, func(q *queries.Queries) error {
		liveroles, err := q.GetLiveRolesByGuildID(ctx, int32(guildId))
		if err != nil {
			return err
		}

		result = lo.Map(liveroles, func(liverole queries.GetLiveRolesByGuildIDRow, _ int) *domain.Liverole {
			return &domain.Liverole{
				ID:            int(liverole.ID),
				GuildID:       guildId,
				RoleID:        int(liverole.RoleID),
				DiscordRoleID: liverole.DiscordRoleID.String,
			}
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l *liverolesRepository) GetLiverole(
	ctx context.Context,
	guildId int,
	discordRoleID string,
) (*domain.Liverole, error) {
	result := &domain.Liverole{}

	err := l.q.Transaction(ctx, func(q *queries.Queries) error {
		liverole, err := q.GetLiveRoleByDiscordID(ctx, queries.GetLiveRoleByDiscordIDParams{
			DiscordRoleID: discordRoleID,
			GuildID:       int32(guildId),
		})
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrLiveroleNotFound
		} else if err != nil {
			return err
		}

		result.ID = int(liverole.ID)
		result.DiscordRoleID = liverole.DiscordRoleID.String
		result.GuildID = int(liverole.GuildID.Int32)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
