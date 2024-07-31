package service

import (
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/repository"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
)

type Service struct {
	repo      repository.Repository
	txManager *txmanager.TxManager
}

func New(
	repo repository.Repository,
	txManager *txmanager.TxManager,
) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
	}
}
