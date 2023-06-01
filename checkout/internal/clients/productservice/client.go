package productservice

import (
	"context"
	"fmt"
	"net/url"
	"route256/checkout/internal/domain"
	"route256/libs/httpclient"
	"time"
)

const (
	listSKUsPath   = "list_skus"
	getProductPath = "get_product"
)

type Client struct {
	urlBase string
	token   string
}

func New(urlBase, token string) *Client {
	return &Client{urlBase: urlBase, token: token}
}

func (c *Client) ListSKUs(ctx context.Context, startAfterSku, Count int64) ([]int64, error) {
	const timeout = 10 * time.Second

	uri, err := url.JoinPath(c.urlBase, listSKUsPath)
	if err != nil {
		return nil, fmt.Errorf("join path: %w", err)
	}

	responseObj := &ListSKUsResponse{}

	err = httpclient.Send(ctx, "POST", uri, timeout, &ListSKUsRequest{
		BaseRequest:   BaseRequest{Token: c.token},
		StartAfterSku: startAfterSku,
		Count:         Count,
	}, responseObj, false)
	if err != nil {
		return nil, fmt.Errorf("product_service: %s: %w", listSKUsPath, err)
	}

	return responseObj.SKUs, nil
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (*domain.ProductSt, error) {
	const timeout = 10 * time.Second

	uri, err := url.JoinPath(c.urlBase, getProductPath)
	if err != nil {
		return nil, fmt.Errorf("join path: %w", err)
	}

	responseObj := &GetProductResponse{}

	err = httpclient.Send(ctx, "POST", uri, timeout, &GetProductRequest{
		BaseRequest: BaseRequest{Token: c.token},
		SKU:         sku,
	}, responseObj, false)
	if err != nil {
		return nil, fmt.Errorf("product_service: %s: %w", getProductPath, err)
	}

	return &domain.ProductSt{
		Name:  responseObj.Name,
		Price: uint32(responseObj.Price),
	}, nil
}
