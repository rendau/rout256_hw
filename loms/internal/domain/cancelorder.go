package domain

import (
	"context"
)

func (m *Model) CancelOrderValidate(orderID int64) error {
	if orderID <= 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (m *Model) CancelOrder(ctx context.Context, orderID int64) error {
	// validate
	if err := m.CancelOrderValidate(orderID); err != nil {
		return err
	}

	return nil
}
