package domain

import (
	"context"
)

func (m *Model) ListOrderValidate(orderID int64) error {
	if orderID <= 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (m *Model) ListOrder(ctx context.Context, orderID int64) (*OrderSt, error) {
	// validate
	if err := m.ListOrderValidate(orderID); err != nil {
		return nil, err
	}

	return &OrderSt{
		Status: "new",
		User:   7,
		Items: []OrderItemSt{
			{Sku: 3, Count: 2},
			{Sku: 4, Count: 1},
		},
	}, nil
}
