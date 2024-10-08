package db

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrorRecordNotFound = pgx.ErrNoRows

var errorCodeNames = map[string]string{
	"23503": "foreign_key_violation",
	"23505": "unique_violation",
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if er, ok := errorCodeNames[pgErr.Code]; ok {
			return er
		}
	}
	return ""
}
