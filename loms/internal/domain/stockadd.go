package domain

import (
	"context"
	"fmt"
)

func (d *Domain) StockAdd(ctx context.Context, warehouseId int64, sku uint32, count uint64) error {
	err := d.repo.StockPut(ctx, warehouseId, sku, count)
	if err != nil {
		return fmt.Errorf("repo.StockPut: %w", err)
	}

	return nil
}
