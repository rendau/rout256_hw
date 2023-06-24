package domain

import (
	"context"

	"route256/checkout/internal/domain/models"
)

//go:generate mockery --name ILomsService --output ./mocks --filename loms_service.go
type ILomsService interface {
	Stocks(ctx context.Context, sku uint32) ([]models.StockSt, error)
	CreateOrder(ctx context.Context, user int64, cart *models.CartSt) (int64, error)
}

//go:generate mockery --name IProductService --output ./mocks --filename product_service.go
type IProductService interface {
	ListSKUs(ctx context.Context, startAfterSku, Count int64) ([]int64, error)
	GetProduct(ctx context.Context, sku int64) (*models.ProductSt, error)
}
