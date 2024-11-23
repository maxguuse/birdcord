package repository

import (
	"context"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TxGetter *txmanager.TxGetter
}

func NewPgx(opts Opts) *tempvoicePgx {
	return &tempvoicePgx{
		txGetter: opts.TxGetter,
	}
}

var _ Repository = &tempvoicePgx{}

type tempvoicePgx struct {
	txGetter *txmanager.TxGetter
}

func (t *tempvoicePgx) GetHubs(ctx context.Context, guildId int64) ([]*domain.TempvoiceHub, error) {
	return []*domain.TempvoiceHub{
		{
			ID:                1,
			DiscordChannelID:  1185121602464141416,
			DiscordGuildID:    1149093125118251018,
			TempvoiceTemplate: "A",
			TempvoiceCategory: 1149093125659312169,
		},
	}, nil
}
