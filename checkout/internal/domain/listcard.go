package domain

import (
	"context"
	"errors"
	"fmt"

	"route256/checkout/internal/domain/models"
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

	// call ProductService to get products
	result.TotalPrice = 0
	for _, item := range result.Items {
		product, err := d.productService.GetProduct(ctx, int64(item.Sku))
		if err != nil {
			return nil, err
		}
		item.Name = product.Name
		item.Price = product.Price
		result.TotalPrice += product.Price * uint32(item.Count)
	}

	return result, nil
}
