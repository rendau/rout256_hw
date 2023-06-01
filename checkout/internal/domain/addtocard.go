package domain

import (
	"context"
	"fmt"
)

func (m *Model) AddToCartValidate(user int64) error {
	if user <= 0 {
		return ErrUserNotFound
	}
	return nil
}

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	// validate
	if err := m.AddToCartValidate(user); err != nil {
		return err
	}

	stocks, err := m.lomsService.Stocks(ctx, sku)
	if err != nil {
		return fmt.Errorf("get stocks: %w", err)
	}

	counter := uint64(count)
	for _, stock := range stocks {
		if counter <= stock.Count {
			return nil
		}
		counter -= stock.Count
	}

	return ErrStockInsufficient
}
