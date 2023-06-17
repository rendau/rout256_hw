package domain

import (
	"context"
	"errors"
	"fmt"

	"route256/checkout/internal/domain/models"
	"route256/libs/workerpool"
)

func (d *Domain) ListCartValidate(user int64) error {
	if user <= 0 {
		return ErrUserNotFound
	}
	return nil
}

func (d *Domain) ListCart(ctx context.Context, user int64) (*models.CartSt, error) {
	// validate
	if err := d.ListCartValidate(user); err != nil {
		return nil, err
	}

	// get cart
	result, err := d.repo.CartGetByUsrId(ctx, user)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			return &models.CartSt{
				User:       user,
				TotalPrice: 0,
				Items:      []*models.CartItemSt{},
			}, nil
		}
		return nil, fmt.Errorf("CartGetByUsrId: %w", err)
	}

	// get cart items
	result.Items, err = d.repo.CartItemList(ctx, models.CartItemListParsSt{CartId: &result.Id})
	if err != nil {
		return nil, fmt.Errorf("CartItemList: %w", err)
	}

	// create worker pool
	wp := workerpool.NewWorkerPool(
		ctx,
		5,
		func(ctx context.Context, cartItem *models.CartItemSt) (*models.ProductSt, error) {
			return d.productService.GetProduct(ctx, int64(cartItem.Sku))
		},
	)

	result.TotalPrice = 0

	go func() {
		for res := range wp.ResultChan() {
			if res.Err != nil {
				continue
			}
			res.Task.Val.Name = res.Val.Name
			res.Task.Val.Price = res.Val.Price
			result.TotalPrice += res.Val.Price * uint32(res.Task.Val.Count)
		}
	}()

	// add tasks
	for _, item := range result.Items {
		wp.AddTask(ctx, item)
	}

	return result, nil
}
