package domain

import (
	"context"
	"route256/libs/logger"

	"route256/checkout/internal/domain/models"
)

func (d *Domain) PurchaseValidate(user int64) error {
	if user <= 0 {
		return ErrUserNotFound
	}
	return nil
}

func (d *Domain) Purchase(ctx context.Context, user int64) (int64, error) {
	// validate
	if err := d.PurchaseValidate(user); err != nil {
		return 0, err
	}

	// get cart
	cart, err := d.repo.CartGetByUsrId(ctx, user)
	if err != nil {
		return 0, err
	}

	// get cart items
	cart.Items, err = d.repo.CartItemList(ctx, models.CartItemListParsSt{
		CartId: &cart.Id,
	})
	if err != nil {
		return 0, err
	}

	// create order
	orderID, err := d.lomsService.CreateOrder(ctx, user, cart)
	if err != nil {
		return 0, err
	}

	// remove cart items
	err = d.repo.CartItemRemoveAllForCartId(ctx, cart.Id)
	if err != nil {
		logger.Errorw(ctx, err, "error removing cart items")
	}

	// remove cart
	err = d.repo.CartRemove(ctx, cart.Id)
	if err != nil {
		logger.Errorw(ctx, err, "error removing cart")
	}

	return orderID, nil
}
