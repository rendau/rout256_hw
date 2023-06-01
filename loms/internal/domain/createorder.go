package domain

import (
	"context"
)

func (m *Model) CreateOrderValidate(resObj *CreateOrderRequestSt) error {
	if resObj.User <= 0 {
		return ErrUserNotFound
	}

	// validate items
	for _, item := range resObj.Items {
		if item.SKU == 0 {
			return ErrSkuRequired
		}
		if item.Count == 0 {
			return ErrCountRequired
		}
	}

	return nil
}

func (m *Model) CreateOrder(ctx context.Context, resObj *CreateOrderRequestSt) (int64, error) {
	// validate
	if err := m.CreateOrderValidate(resObj); err != nil {
		return 0, err
	}

	return 0, nil
}
