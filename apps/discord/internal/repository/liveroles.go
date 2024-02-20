package repository

import (
	"context"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
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
	DeleteLiverole(
		ctx context.Context,
		guildID, roleID int,
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
	panic("TODO: Implement")
}

func (l *liverolesRepository) DeleteLiverole(
	ctx context.Context,
	guildID int,
	roleID int,
) error {
	panic("TODO: Implement")
}

func (l *liverolesRepository) DeleteLiveroles(
	ctx context.Context,
	guildID int,
	discordRolesIds []string,
) error {
	panic("TODO: Implement")
}

func (l *liverolesRepository) GetLiveroles(
	ctx context.Context,
	guildId int,
) ([]*domain.Liverole, error) {
	panic("TODO: Implement")
}
