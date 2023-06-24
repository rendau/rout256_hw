package domain

import (
	"context"
	"errors"
	"reflect"
	domainMocks "route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/domain/models"
	repoMock "route256/checkout/internal/repo/mock"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestDomain_ListCart(t *testing.T) {
	userId := int64(1)
	cart1Id := int64(1)

	ctx := context.Background()

	someErr := errors.New("some error")

	tests := []struct {
		cart         *models.CartSt
		cartErr      error
		cartItems    []*models.CartItemSt
		cartItemsErr error
		products     map[int64]*models.ProductSt
		productsErr  map[int64]error

		argUser int64
		want    *models.CartSt
		wantErr error
	}{
		{
			cart:      &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{},
			products:  map[int64]*models.ProductSt{},

			argUser: -5,
			want:    nil,
			wantErr: ErrUserNotFound,
		},
		{
			cartErr: ErrCartNotFound,

			argUser: userId,
			want: &models.CartSt{
				User:       userId,
				TotalPrice: 0,
				Items:      []*models.CartItemSt{},
			},
		},
		{
			cartErr: someErr,

			argUser: userId,
			wantErr: someErr,
		},
		{
			cart:         &models.CartSt{Id: cart1Id, User: userId},
			cartItemsErr: someErr,

			argUser: userId,
			wantErr: someErr,
		},
		{
			cart: &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{
				{CartId: cart1Id, Sku: 1, Count: 1},
				{CartId: cart1Id, Sku: 2, Count: 2},
				{CartId: cart1Id, Sku: 3, Count: 3},
				{CartId: cart1Id, Sku: 4, Count: 4},
				{CartId: cart1Id, Sku: 5, Count: 5},
				{CartId: cart1Id, Sku: 6, Count: 6},
				{CartId: cart1Id, Sku: 7, Count: 7},
				{CartId: cart1Id, Sku: 8, Count: 8},
				{CartId: cart1Id, Sku: 9, Count: 9},
				{CartId: cart1Id, Sku: 10, Count: 10},
			},
			products: map[int64]*models.ProductSt{
				1:  {Name: "product1", Price: 100},
				2:  {Name: "product2", Price: 50},
				3:  {Name: "product3", Price: 100},
				4:  {Name: "product4", Price: 50},
				5:  {Name: "product5", Price: 100},
				6:  {Name: "product6", Price: 50},
				7:  {Name: "product7", Price: 100},
				8:  {Name: "product8", Price: 50},
				9:  {Name: "product9", Price: 100},
				10: {Name: "product10", Price: 50},
			},
			productsErr: map[int64]error{
				2: someErr,
			},

			argUser: userId,
			wantErr: someErr,
		},
		{
			cart: &models.CartSt{Id: cart1Id, User: userId},
			cartItems: []*models.CartItemSt{
				{CartId: cart1Id, Sku: 1, Count: 1},
				{CartId: cart1Id, Sku: 2, Count: 2},
			},
			products: map[int64]*models.ProductSt{
				1: {Name: "product1", Price: 100},
				2: {Name: "product2", Price: 50},
			},

			argUser: userId,
			want: &models.CartSt{
				Id:         cart1Id,
				User:       userId,
				TotalPrice: 200,
				Items: []*models.CartItemSt{
					{
						CartId: cart1Id,
						Sku:    1,
						Count:  1,
						Name:   "product1",
						Price:  100,
					},
					{
						CartId: cart1Id,
						Sku:    2,
						Count:  2,
						Name:   "product2",
						Price:  50,
					},
				},
			},
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
			for itemId, pr := range tt.products {
				productService.On("GetProduct", mock.Anything, itemId).Return(pr, tt.productsErr[itemId]).After(50 * time.Millisecond)
			}

			got, err := domain.ListCart(ctx, tt.argUser)
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ListCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListCart() got = %v, want %v", got, tt.want)
			}
		})
	}
}
