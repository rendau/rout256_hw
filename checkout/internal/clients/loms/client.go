package loms

import (
	"context"
	"fmt"
	"net/url"
	"route256/checkout/internal/domain"
	"route256/libs/httpclient"
	"time"
)

const (
	stocksPath      = "stocks"
	createOrderPath = "createOrder"
)

type Client struct {
	urlBase string
}

func New(urlBase string) *Client {
	return &Client{urlBase: urlBase}
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.StockSt, error) {
	const timeout = 5 * time.Second

	uri, err := url.JoinPath(c.urlBase, stocksPath)
	if err != nil {
		return nil, fmt.Errorf("join path: %w", err)
	}

	responseObj := &StocksResponse{}

	err = httpclient.Send(ctx, "POST", uri, timeout, &StocksRequest{SKU: sku}, responseObj, true)
	if err != nil {
		return nil, fmt.Errorf("loms: %s: %w", stocksPath, err)
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

func (c *Client) CreateOrder(ctx context.Context, user int64, cart *domain.CartSt) (int64, error) {
	const timeout = 5 * time.Second

	uri, err := url.JoinPath(c.urlBase, createOrderPath)
	if err != nil {
		return 0, fmt.Errorf("join path: %w", err)
	}

	requestObj := &CreateOrderRequest{
		User:  user,
		Items: make([]CreateOrderRequestItem, len(cart.Items)),
	}

	for i, v := range cart.Items {
		requestObj.Items[i].SKU = v.SKU
		requestObj.Items[i].Count = v.Count
	}

	responseObj := &CreateOrderResponse{}

	err = httpclient.Send(ctx, "POST", uri, timeout, requestObj, responseObj, true)
	if err != nil {
		return 0, fmt.Errorf("loms: %s: %w", createOrderPath, err)
	}

	return responseObj.OrderID, nil
}
