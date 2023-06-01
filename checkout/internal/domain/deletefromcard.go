package domain

import (
	"context"
)

func (m *Model) DeleteFromCartValidate(user int64, sku uint32, count uint16) error {
	if user == 0 {
		return ErrUserNotFound
	}
	if sku == 0 {
		return ErrSkuRequired
	}
	if count == 0 {
		return ErrCountRequired
	}
	return nil
}

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	// validate
	if err := m.DeleteFromCartValidate(user, sku, count); err != nil {
		return err
	}

	return nil
}
