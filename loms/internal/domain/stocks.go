package domain

import (
	"context"
)

func (m *Model) StocksValidate(sku uint32) error {
	if sku == 0 {
		return ErrSkuRequired
	}
	return nil
}

func (m *Model) Stocks(ctx context.Context, sku uint32) (*StocksResponseSt, error) {
	// validate
	if err := m.StocksValidate(sku); err != nil {
		return nil, err
	}

	return &StocksResponseSt{
		Stocks: []StockSt{
			{WarehouseID: 1, Count: 100},
		},
	}, nil
}
