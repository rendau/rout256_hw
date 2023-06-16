package domain

import (
	"context"
	"fmt"
)

func (d *Domain) StockRemove(ctx context.Context, warehouseId int64, sku uint32) error {
	err := d.repo.StockRemove(ctx, warehouseId, sku)
	if err != nil {
		return fmt.Errorf("repo.StockRemove: %w", err)
	}

	return nil
}
