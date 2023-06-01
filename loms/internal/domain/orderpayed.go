package domain

import (
	"context"
)

func (m *Model) OrderPayedValidate(orderID int64) error {
	if orderID <= 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (m *Model) OrderPayed(ctx context.Context, orderID int64) error {
	// validate
	if err := m.OrderPayedValidate(orderID); err != nil {
		return err
	}

	return nil
}
