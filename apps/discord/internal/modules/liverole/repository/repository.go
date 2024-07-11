package repository

import (
	"context"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

type Repository interface {
	CreateLiverole(
		ctx context.Context,
		discordGuildId string,
		discordRoleId string,
	) error
	GetLiveroles(
		ctx context.Context,
		discordGuildId string,
	) ([]*domain.Liverole, error)
	GetLiverole(
		ctx context.Context,
		discordGuildId string,
		discordRoleId string,
	) (*domain.Liverole, error)
	DeleteLiveroles(
		ctx context.Context,
		discordGuildId string,
		discordRolesIds []string,
	) error
}
