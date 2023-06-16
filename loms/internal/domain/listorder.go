package domain

import (
	"context"

	"route256/loms/internal/domain/models"
)

func (d *Domain) ListOrderValidate(orderID int64) error {
	if orderID <= 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (d *Domain) ListOrder(ctx context.Context, orderID int64) (*models.OrderSt, error) {
	// validate
	if err := d.ListOrderValidate(orderID); err != nil {
		return nil, err
	}

	// get order
	order, err := d.repo.OrderGet(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// get order items
	order.Items, err = d.repo.OrderGetItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
