package loms

import (
	"context"
	"fmt"

	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/models"
	"route256/checkout/pkg/proto/loms_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client loms_v1.LomsClient
}

func New(uri string) (*Client, error) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return &Client{
		client: loms_v1.NewLomsClient(conn),
	}, nil
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.StockSt, error) {
	responseObj, err := c.client.Stocks(ctx, &loms_v1.StocksRequest{Sku: sku})
	if err != nil {
		return nil, fmt.Errorf("loms: Stocks: %w", err)
	}

	result := make([]domain.StockSt, 0, len(responseObj.Stocks))
	for _, v := range responseObj.Stocks {
		result = append(result, domain.StockSt{
			WarehouseID: v.WarehouseID,
			Count:       v.Count,
		})
	}

	return result, nil
}

func (c *Client) CreateOrder(ctx context.Context, user int64, cart *models.CartSt) (int64, error) {
	requestObj := &loms_v1.CreateOrderRequest{
		User:  user,
		Items: make([]*loms_v1.Order, len(cart.Items)),
	}

	for i, v := range cart.Items {
		requestObj.Items[i] = &loms_v1.Order{
			Sku:   v.Sku,
			Count: uint32(v.Count),
		}
	}

	responseObj, err := c.client.CreateOrder(ctx, requestObj)
	if err != nil {
		return 0, fmt.Errorf("loms: CreateOrder: %w", err)
	}

	return responseObj.OrderID, nil
}
