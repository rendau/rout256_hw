package domain

import (
	"route256/checkout/internal/repo"
	"route256/libs/db"
)

type Domain struct {
	db             db.Transaction
	repo           repo.Repo
	lomsService    ILomsService
	productService IProductService
}

func New(
	db db.Transaction,
	repo repo.Repo,
	stockChecker ILomsService,
	productService IProductService,
) *Domain {
	return &Domain{
		db:             db,
		repo:           repo,
		lomsService:    stockChecker,
		productService: productService,
	}
}
