package domain

import (
	"context"
)

func (m *Model) PurchaseValidate(user int64) error {
	if user <= 0 {
		return ErrUserNotFound
	}
	return nil
}

func (m *Model) Purchase(ctx context.Context, user int64) (int64, error) {
	// validate
	if err := m.PurchaseValidate(user); err != nil {
		return 0, err
	}

	// get cart
	cart, err := m.ListCart(ctx, user)
	if err != nil {
		return 0, err
	}

	// create order
	orderID, err := m.lomsService.CreateOrder(ctx, user, cart)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
