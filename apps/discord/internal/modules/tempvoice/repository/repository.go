package repository

import (
	"context"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

type Repository interface {
	GetHubs(context.Context, int64) ([]*domain.TempvoiceHub, error)
}
