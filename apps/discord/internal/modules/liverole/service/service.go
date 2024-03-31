package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/samber/lo"
)

type Service struct {
	db repository.DB
}

func New(db repository.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Add(ctx context.Context, r *AddLiveRoleRequest) error {
	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, r.GuildID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	_, err = s.db.Liveroles().CreateLiverole(ctx, r.RoleID, guild.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrRoleAlreadyExists
		}

		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *Service) Clear(ctx context.Context, guildID string) error {
	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, guildID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	liveroles, err := s.db.Liveroles().GetLiveroles(ctx, guild.ID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	if len(liveroles) == 0 {
		return ErrNoLiveroles
	}

	err = s.db.Liveroles().DeleteLiveroles(
		ctx,
		guild.ID,
		lo.Map(liveroles, func(liverole *domain.Liverole, _ int) string {
			return liverole.DiscordRoleID
		}))
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, guildID string) ([]string, error) {
	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, guildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	liveroles, err := s.db.Liveroles().GetLiveroles(ctx, guild.ID)
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
	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, r.GuildID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	role, err := s.db.Liveroles().GetLiverole(ctx, guild.ID, r.RoleID)
	if errors.Is(err, repository.ErrLiveroleNotFound) {
		return ErrLiveroleNotFound
	}

	err = s.db.Liveroles().DeleteLiverole(ctx, role.ID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
