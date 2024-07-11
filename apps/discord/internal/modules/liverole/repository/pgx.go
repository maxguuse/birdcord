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
	discordGuildId string,
	discordRoleId string,
) error {
	err := l.txm.Do(ctx, func(db qrm.DB) error {
		_, err := Liveroles.INSERT(
			Liveroles.RoleID,
		).VALUES(
			postgres.SELECT(
				Roles.ID,
			).FROM(
				Roles.LEFT_JOIN(
					Guilds,
					Roles.GuildID.EQ(Guilds.ID),
				),
			).WHERE(
				Guilds.DiscordGuildID.
					EQ(postgres.String(discordGuildId)).
					AND(Roles.DiscordRoleID.
						EQ(postgres.String(discordRoleId))),
			),
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
	Liveroles.ID.AS("Liverole.Id"),
	Guilds.ID.AS("Liverole.GuildId"),
	Roles.ID.AS("Liverole.RoleId"),
	Roles.DiscordRoleID.AS("Liverole.DiscordRoleId"),
).FROM(Liveroles.
	LEFT_JOIN(Roles, Liveroles.RoleID.EQ(Roles.ID)).
	LEFT_JOIN(Guilds, Roles.GuildID.EQ(Guilds.ID)),
)

func (l *liverolePgx) GetLiveroles(
	ctx context.Context,
	discordGuildId string,
) ([]*domain.Liverole, error) {
	var liveroles []*domain.Liverole
	err := l.txm.Do(ctx, func(db qrm.DB) error {
		err := liveroleSelect.WHERE(
			Guilds.DiscordGuildID.EQ(postgres.String(discordGuildId)),
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
	discordGuildId string,
	discordRoleId string,
) (*domain.Liverole, error) {
	var liverole *domain.Liverole
	err := l.txm.Do(ctx, func(db qrm.DB) error {
		err := liveroleSelect.WHERE(
			Guilds.DiscordGuildID.EQ(
				postgres.String(discordGuildId),
			).AND(Roles.DiscordRoleID.EQ(
				postgres.String(discordRoleId)),
			)).QueryContext(ctx, db, liverole)
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
	discordGuildId string,
	discordRolesIds []string,
) error {
	rolesExpr := lo.Map(discordRolesIds, func(discordRoleId string, _ int) postgres.Expression {
		return postgres.String(discordRoleId)
	})

	err := l.txm.Do(ctx, func(db qrm.DB) error {
		_, err := Liveroles.DELETE().USING(
			Roles,
			Guilds,
		).WHERE(
			Liveroles.RoleID.EQ(Roles.ID).
				AND(Roles.DiscordRoleID.IN(rolesExpr...)).
				AND(Roles.GuildID.EQ(Guilds.ID)).
				AND(Guilds.DiscordGuildID.EQ(postgres.String(discordGuildId))),
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
