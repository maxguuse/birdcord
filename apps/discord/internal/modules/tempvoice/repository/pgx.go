package repository

import (
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
