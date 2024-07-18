package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/repository"
	"github.com/maxguuse/birdcord/libs/jet/pgerrors"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/samber/lo"
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

func (s *Service) Add(ctx context.Context, r *AddLiveRoleRequest) error {
	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	roleId, err := strconv.Atoi(r.RoleID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.repo.CreateLiverole(
		ctx,
		int64(guildId),
		int64(roleId),
	)

	if errors.Is(err, pgerrors.ErrDuplicateKey) {
		return ErrRoleAlreadyExists
	}

	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *Service) Clear(ctx context.Context, discordGuildId string) error {
	guildId, err := strconv.Atoi(discordGuildId)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		liveroles, err := s.repo.GetLiveroles(ctx, int64(guildId))
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		if len(liveroles) == 0 {
			return ErrNoLiveroles
		}

		err = s.repo.DeleteLiveroles(
			ctx,
			int64(guildId),
			lo.Map(liveroles, func(liverole *domain.Liverole, _ int) int64 {
				return int64(liverole.DiscordRoleID)
			}),
		)
		if err != nil {
			return errors.Join(domain.ErrInternal, err)
		}

		return nil
	})

	if errors.Is(err, txmanager.ErrTx) {
		return errors.Join(domain.ErrInternal, err)
	}

	return err
}

func (s *Service) List(ctx context.Context, discordGuildId string) ([]string, error) {
	guildId, err := strconv.Atoi(discordGuildId)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	liveroles, err := s.repo.GetLiveroles(ctx, int64(guildId))
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	if len(liveroles) == 0 {
		return nil, ErrNoLiveroles
	}

	rolesList := lo.Map(liveroles, func(liverole *domain.Liverole, _ int) string {
		return fmt.Sprintf("<@&%d>", liverole.DiscordRoleID)
	})

	return rolesList, nil
}

func (s *Service) Remove(ctx context.Context, r *RemoveLiveRoleRequest) error {
	guildId, err := strconv.Atoi(r.GuildID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	roleId, err := strconv.Atoi(r.RoleID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.repo.DeleteLiveroles(ctx, int64(guildId), []int64{int64(roleId)})
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
