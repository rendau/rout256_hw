package domain

import (
	"context"

	"route256/libs/constant"
)

func (d *Domain) OrderPayedValidate(orderID int64) error {
	if orderID <= 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (d *Domain) OrderPayed(ctx context.Context, orderID int64) error {
	var err error

	// validate
	if err = d.OrderPayedValidate(orderID); err != nil {
		return err
	}

	err = d.repo.OrderSetStatus(ctx, orderID, constant.OrderStatusPayed)
	if err != nil {
		return err
	}

	return nil
}
