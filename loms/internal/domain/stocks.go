package domain

import (
	"context"
	"fmt"

	"route256/loms/internal/domain/models"
)

func (d *Domain) StocksValidate(sku uint32) error {
	if sku == 0 {
		return ErrSkuRequired
	}
	return nil
}

func (d *Domain) Stocks(ctx context.Context, sku uint32) ([]*models.StockSt, error) {
	// validate
	if err := d.StocksValidate(sku); err != nil {
		return nil, err
	}

	items, err := d.repo.StockList(ctx, models.StockListParsSt{Sku: &sku}, false)
	if err != nil {
		return nil, fmt.Errorf("repo.StockList: %w", err)
	}

	return items, nil
}
