package domain

import (
	"context"
	"errors"
	domainMocks "route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/domain/models"
	repoMock "route256/checkout/internal/repo/mock"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestDomain_Purchase(t *testing.T) {
	userId := int64(1)
	cart1Id := int64(1)

	ctx := context.Background()

	someErr := errors.New("some error")

	tests := []struct {
		cart           *models.CartSt
		cartErr        error
		cartItems      []*models.CartItemSt
		cartItemsErr   error
		orderId        int64
		orderErr       error
		itemsRemoveErr error
		cartRemoveErr  error

		argUser int64
		want    int64
		wantErr error
	}{
		{
			argUser: -1,
			wantErr: ErrUserNotFound,
		},
		{
			cartErr: someErr,

			argUser: 1,
			wantErr: someErr,
		},
		{
			cart:         &models.CartSt{Id: cart1Id, User: userId},
			cartItemsErr: someErr,

			argUser: 1,
			wantErr: someErr,
		},
		{
			cart: &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{
				{CartId: cart1Id, Sku: 1, Count: 1},
			},
			orderErr: someErr,

			argUser: 1,
			wantErr: someErr,
		},
		{
			cart: &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{
				{CartId: cart1Id, Sku: 1, Count: 1},
			},
			orderId:        1,
			itemsRemoveErr: someErr,

			argUser: 1,
			want:    1,
			wantErr: nil,
		},
		{
			cart: &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{
				{CartId: cart1Id, Sku: 1, Count: 1},
			},
			orderId:       1,
			cartRemoveErr: someErr,

			argUser: 1,
			want:    1,
			wantErr: nil,
		},
		{
			cart: &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{
				{CartId: cart1Id, Sku: 1, Count: 1},
				{CartId: cart1Id, Sku: 2, Count: 2},
			},
			orderId: 1,

			argUser: 1,
			want:    1,
			wantErr: nil,
		},
	}
	for ttI, tt := range tests {
		t.Run(strconv.Itoa(ttI+1), func(t *testing.T) {
			repo := &repoMock.Repo{}
			lomsService := &domainMocks.ILomsService{}
			productService := &domainMocks.IProductService{}
			domain := New(repo, lomsService, productService)

			// set mocks
			repo.On("CartGetByUsrId", mock.Anything, userId).Return(tt.cart, tt.cartErr)
			repo.On("CartItemList", mock.Anything, models.CartItemListParsSt{CartId: &cart1Id}).Return(tt.cartItems, tt.cartItemsErr)
			lomsService.On("CreateOrder", ctx, userId, mock.Anything).Return(tt.orderId, tt.orderErr)
			repo.On("CartItemRemoveAllForCartId", mock.Anything, cart1Id).Return(tt.itemsRemoveErr)
			repo.On("CartRemove", mock.Anything, cart1Id).Return(tt.cartRemoveErr)

			got, err := domain.Purchase(ctx, tt.argUser)
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Purchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Purchase() got = %v, want %v", got, tt.want)
			}

			if tt.cartErr == nil && tt.cartItemsErr == nil && tt.wantErr != ErrUserNotFound {
				lomsService.AssertCalled(t, "CreateOrder", ctx, userId, mock.MatchedBy(func(cart *models.CartSt) bool {
					if !(cart.Id == cart1Id && cart.User == userId) {
						return false
					}
					if len(cart.Items) != len(tt.cartItems) {
						return false
					}
					for i, item := range tt.cartItems {
						if *item != *cart.Items[i] {
							return false
						}
					}
					return true
				}))
			}
		})
	}
}
