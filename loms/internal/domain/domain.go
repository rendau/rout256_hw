package domain

import (
	"route256/libs/db"
	"route256/loms/internal/repo"
)

type Domain struct {
	db   db.Transaction
	repo repo.Repo
}

func New(db db.Transaction, repo repo.Repo) *Domain {
	return &Domain{
		db:   db,
		repo: repo,
	}
}
