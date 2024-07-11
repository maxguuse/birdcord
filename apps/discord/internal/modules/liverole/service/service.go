package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/repository"
	db "github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/samber/lo"
)

type Service struct {
	repo repository.Repository
}

func New(
	repo repository.Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, r *AddLiveRoleRequest) error {
	err := s.repo.CreateLiverole(ctx, r.GuildID, r.RoleID)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return ErrRoleAlreadyExists // TODO: Add this error to repository
		}

		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *Service) Clear(ctx context.Context, guildID string) error {
	liveroles, err := s.repo.GetLiveroles(ctx, guildID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	if len(liveroles) == 0 {
		return ErrNoLiveroles
	}

	err = s.repo.DeleteLiveroles(
		ctx,
		guildID,
		lo.Map(liveroles, func(liverole *domain.Liverole, _ int) string {
			return liverole.DiscordRoleID
		}),
	)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, guildID string) ([]string, error) {
	liveroles, err := s.repo.GetLiveroles(ctx, guildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if len(liveroles) == 0 {
		return nil, ErrNoLiveroles
	}

	rolesList := lo.Map(liveroles, func(liverole *domain.Liverole, _ int) string {
		return fmt.Sprintf("<@&%s>", liverole.DiscordRoleID)
	})

	return rolesList, nil
}

func (s *Service) Remove(ctx context.Context, r *RemoveLiveRoleRequest) error {
	err := s.repo.DeleteLiveroles(ctx, r.GuildID, []string{r.RoleID})
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
