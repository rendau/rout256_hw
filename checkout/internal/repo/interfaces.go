package repo

import (
	"context"

	"route256/checkout/internal/domain/models"
)

//go:generate mockery --name Repo --output ./mock --filename repo.go
type Repo interface {
	// Cart
	CartList(ctx context.Context, pars models.CartListParsSt) ([]*models.CartSt, error)
	CartGet(ctx context.Context, id int64) (*models.CartSt, error)
	CartGetByUsrId(ctx context.Context, userId int64) (*models.CartSt, error)
	CartCreate(ctx context.Context, obj *models.CartSt) (int64, error)
	CartRemove(ctx context.Context, id int64) error

	// CartItem
	CartItemList(ctx context.Context, pars models.CartItemListParsSt) ([]*models.CartItemSt, error)
	CartItemGet(ctx context.Context, cartId int64, sku uint32) (*models.CartItemSt, error)
	CartItemSet(ctx context.Context, obj *models.CartItemSt) error
	CartItemRemove(ctx context.Context, cartId int64, sku uint32) error
	CartItemRemoveAllForCartId(ctx context.Context, cartId int64) error
}
