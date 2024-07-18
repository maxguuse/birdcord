package repository

import (
	"context"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

type Repository interface {
	CreateLiverole(
		ctx context.Context,
		discordGuildId, discordRoleId int64,
	) error
	GetLiveroles(
		ctx context.Context,
		discordGuildId int64,
	) ([]*domain.Liverole, error)
	GetLiverole(
		ctx context.Context,
		discordGuildId, discordRoleId int64,
	) (*domain.Liverole, error)
	DeleteLiveroles(
		ctx context.Context,
		discordGuildId int64,
		discordRolesIds []int64,
	) error
}
