package pgerrors

import (
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var Codes = map[string]error{
	"23505": ErrDuplicateKey,
}

var (
	ErrDuplicateKey = errors.New("duplicate key")
	ErrNotFound     = errors.New("resource not found")
)

func PgErrorOrErr(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		pgerr, ok := Codes[pgErr.Code]
		if ok {
			return pgerr
		}
	}

	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) || errors.Is(err, qrm.ErrNoRows) {
		return ErrNotFound
	}

	return err
}
