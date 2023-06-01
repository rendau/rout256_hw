package domain

import (
	"context"
)

func (m *Model) ListCartValidate(user int64) error {
	if user <= 0 {
		return ErrUserNotFound
	}
	return nil
}

func (m *Model) ListCart(ctx context.Context, user int64) (*CartSt, error) {
	// validate
	if err := m.ListCartValidate(user); err != nil {
		return nil, err
	}

	result := &CartSt{
		Items: []*CartItemSt{
			{SKU: 4678816, Count: 1},
			{SKU: 4288068, Count: 2},
			{SKU: 4487693, Count: 3},
		},
		TotalPrice: 10000,
	}

	// call ProductService to get products
	for _, item := range result.Items {
		product, err := m.productService.GetProduct(ctx, int64(item.SKU))
		if err != nil {
			return nil, err
		}
		item.Name = product.Name
		item.Price = product.Price
	}

	return result, nil
}
