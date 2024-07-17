package repository

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	. "github.com/maxguuse/birdcord/libs/jet/generated/birdcord/public/table"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TxManager *txmanager.TxManager
}

func NewPgx(opts Opts) *liverolePgx {
	return &liverolePgx{
		txm: opts.TxManager,
	}
}

var _ Repository = &liverolePgx{}

type liverolePgx struct {
	txm *txmanager.TxManager
}

func (l *liverolePgx) CreateLiverole(
	ctx context.Context,
	discordGuildId int64,
	discordRoleId int64,
) error {
	err := l.txm.Do(ctx, func(db qrm.DB) error {
		_, err := Liveroles.INSERT(
			Liveroles.DiscordRoleID,
			Liveroles.DiscordGuildID,
		).VALUES(
			discordRoleId,
			discordGuildId,
		).ExecContext(ctx, db)
		if err != nil {
			return err // TODO: Wrap error
		}

		return nil
	})
	if err != nil {
		return err // TODO: Wrap error
	}

	return nil
}

var liveroleSelect = postgres.SELECT(
	Liveroles.ID,
	Liveroles.DiscordRoleID,
	Liveroles.DiscordGuildID,
).FROM(
	Liveroles,
)

func (l *liverolePgx) GetLiveroles(
	ctx context.Context,
	discordGuildId int64,
) ([]*domain.Liverole, error) {
	var liveroles []*domain.Liverole
	err := l.txm.Do(ctx, func(db qrm.DB) error {
		err := liveroleSelect.WHERE(
			Liveroles.DiscordGuildID.EQ(postgres.Int64(discordGuildId)),
		).QueryContext(ctx, db, &liveroles)
		if err != nil {
			return err // TODO: Wrap error
		}

		return nil
	})
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return liveroles, nil
}

func (l *liverolePgx) GetLiverole(
	ctx context.Context,
	discordGuildId int64,
	discordRoleId int64,
) (*domain.Liverole, error) {
	var liverole *domain.Liverole
	err := l.txm.Do(ctx, func(db qrm.DB) error {
		err := liveroleSelect.WHERE(
			Liveroles.DiscordGuildID.EQ(postgres.Int64(discordGuildId)).
				AND(Liveroles.DiscordRoleID.EQ(postgres.Int64(discordRoleId))),
		).QueryContext(ctx, db, liverole)
		if err != nil {
			return err // TODO: Wrap error
		}

		return nil
	})
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return liverole, nil
}

func (l *liverolePgx) DeleteLiveroles(
	ctx context.Context,
	discordGuildId int64,
	discordRolesIds []int64,
) error {
	rolesExpr := lo.Map(discordRolesIds, func(discordRoleId int64, _ int) postgres.Expression {
		return postgres.Int64(discordRoleId)
	})

	err := l.txm.Do(ctx, func(db qrm.DB) error {
		_, err := Liveroles.DELETE().WHERE(
			Liveroles.DiscordRoleID.IN(rolesExpr...).
				AND(Liveroles.DiscordGuildID.EQ(postgres.Int64(discordGuildId))),
		).ExecContext(ctx, db)
		if err != nil {
			return err // TODO: Wrap error
		}

		return nil
	})
	if err != nil {
		return err // TODO: Wrap error
	}

	return nil
}
