package repository

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	. "github.com/maxguuse/birdcord/libs/jet/generated/birdcord/public/table"
	"github.com/maxguuse/birdcord/libs/jet/pgerrors"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TxGetter *txmanager.TxGetter
}

func NewPgx(opts Opts) *liverolePgx {
	return &liverolePgx{
		txGetter: opts.TxGetter,
	}
}

var _ Repository = &liverolePgx{}

type liverolePgx struct {
	txGetter *txmanager.TxGetter
}

func (l *liverolePgx) CreateLiverole(
	ctx context.Context,
	discordGuildId int64,
	discordRoleId int64,
) error {
	_, err := Liveroles.INSERT(
		Liveroles.DiscordRoleID,
		Liveroles.DiscordGuildID,
	).VALUES(
		discordRoleId,
		discordGuildId,
	).ExecContext(ctx, l.txGetter.DefaultTxOrDB(ctx))

	if pgerr := pgerrors.PgErrorOrErr(err); pgerr != nil {
		return pgerr
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
	err := liveroleSelect.WHERE(
		Liveroles.DiscordGuildID.EQ(postgres.Int64(discordGuildId)),
	).QueryContext(ctx, l.txGetter.DefaultTxOrDB(ctx), &liveroles)

	if pgerr := pgerrors.PgErrorOrErr(err); pgerr != nil {
		return nil, pgerr
	}

	return liveroles, nil
}

func (l *liverolePgx) GetLiverole(
	ctx context.Context,
	discordGuildId int64,
	discordRoleId int64,
) (*domain.Liverole, error) {
	var liverole *domain.Liverole
	err := liveroleSelect.WHERE(
		Liveroles.DiscordGuildID.EQ(postgres.Int64(discordGuildId)).
			AND(Liveroles.DiscordRoleID.EQ(postgres.Int64(discordRoleId))),
	).QueryContext(ctx, l.txGetter.DefaultTxOrDB(ctx), liverole)

	if pgerr := pgerrors.PgErrorOrErr(err); pgerr != nil {
		return nil, pgerr
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

	_, err := Liveroles.DELETE().WHERE(
		Liveroles.DiscordRoleID.IN(rolesExpr...).
			AND(Liveroles.DiscordGuildID.EQ(postgres.Int64(discordGuildId))),
	).ExecContext(ctx, l.txGetter.DefaultTxOrDB(ctx))

	if pgerr := pgerrors.PgErrorOrErr(err); pgerr != nil {
		return pgerr
	}

	return nil
}
