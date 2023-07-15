package pg

import (
	"github.com/Masterminds/squirrel"

	"route256/libs/db"
)

type St struct {
	db db.DB
	sq squirrel.StatementBuilderType
}

func New(db db.DB) *St {
	return &St{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
