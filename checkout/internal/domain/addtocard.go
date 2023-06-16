package domain

import (
	"context"
	"errors"
	"fmt"

	"route256/checkout/internal/domain/models"
)

func (d *Domain) AddToCartValidate(user int64) error {
	if user <= 0 {
		return ErrUserNotFound
	}
	return nil
}

func (d *Domain) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	// validate
	if err := d.AddToCartValidate(user); err != nil {
		return err
	}

	stocks, err := d.lomsService.Stocks(ctx, sku)
	if err != nil {
		return fmt.Errorf("get stocks: %w", err)
	}

	counter := uint64(count)
	for _, stock := range stocks {
		if counter <= stock.Count {
			counter = 0
			break
		}
		counter -= stock.Count
	}
	if counter > 0 {
		return ErrStockInsufficient
	}

	var cartId int64

	// get cart
	cart, err := d.repo.CartGetByUsrId(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, ErrCartNotFound):
			cartId, err = d.repo.CartCreate(ctx, &models.CartSt{
				User: user,
			})
			if err != nil {
				return fmt.Errorf("create cart: %w", err)
			}
		default:
			return fmt.Errorf("get cart: %w", err)
		}
	} else {
		cartId = cart.Id
	}

	// get cart item
	cartItem, err := d.repo.CartItemGet(ctx, cartId, sku)
	if err != nil {
		return fmt.Errorf("get cart item: %w", err)
	}
	if cartItem != nil {
		count += cartItem.Count
	}

	err = d.repo.CartItemSet(ctx, &models.CartItemSt{
		CartId: cartId,
		Sku:    sku,
		Count:  count,
	})
	if err != nil {
		return fmt.Errorf("add cart item: %w", err)
	}

	return nil
}
