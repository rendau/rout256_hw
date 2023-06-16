package domain

import (
	"context"
	"errors"
	"fmt"

	"route256/checkout/internal/domain/models"
)

func (d *Domain) DeleteFromCartValidate(user int64, sku uint32, count uint16) error {
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

func (d *Domain) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	// validate
	if err := d.DeleteFromCartValidate(user, sku, count); err != nil {
		return err
	}

	// get cart
	cart, err := d.repo.CartGetByUsrId(ctx, user)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			return nil
		}
		return fmt.Errorf("get cart: %w", err)
	}

	// get cart item
	cartItem, err := d.repo.CartItemGet(ctx, cart.Id, sku)
	if err != nil {
		return fmt.Errorf("get cart item: %w", err)
	}
	if cartItem == nil {
		return nil
	}

	if count >= cartItem.Count {
		err = d.repo.CartItemRemove(ctx, cart.Id, sku)
		if err != nil {
			return fmt.Errorf("remove cart item: %w", err)
		}
	} else {
		err = d.repo.CartItemSet(ctx, &models.CartItemSt{
			CartId: cart.Id,
			Sku:    sku,
			Count:  cartItem.Count - count,
		})
		if err != nil {
			return fmt.Errorf("set cart item: %w", err)
		}
	}

	return nil
}
