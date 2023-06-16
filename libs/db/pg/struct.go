package pg

import (
	"github.com/jackc/pgx/v4"
)

type rowsSt struct {
	pgx.Rows
}

func (o *rowsSt) Err() error {
	return o.Rows.Err()
}

func (o *rowsSt) Scan(dest ...any) error {
	return o.Rows.Scan(dest...)
}

type rowSt struct {
	pgx.Row
}

func (o *rowSt) Scan(dest ...any) error {
	return o.Row.Scan(dest...)
}
