package domain

import (
	"route256/checkout/internal/repo"
)

type Domain struct {
	repo           repo.Repo
	lomsService    ILomsService
	productService IProductService
}

func New(
	repo repo.Repo,
	stockChecker ILomsService,
	productService IProductService,
) *Domain {
	return &Domain{
		repo:           repo,
		lomsService:    stockChecker,
		productService: productService,
	}
}
