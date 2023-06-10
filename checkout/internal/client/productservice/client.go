package productservice

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/proto/product_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client product_service.ProductServiceClient
	token  string
}

func New(uri, token string) (*Client, error) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return &Client{
		client: product_service.NewProductServiceClient(conn),
		token:  token,
	}, nil
}

func (c *Client) ListSKUs(ctx context.Context, startAfterSku, Count int64) ([]int64, error) {
	responseObj, err := c.client.ListSkus(ctx, &product_service.ListSkusRequest{
		Token:         c.token,
		StartAfterSku: uint32(startAfterSku),
		Count:         uint32(Count),
	})
	if err != nil {
		return nil, fmt.Errorf("product_service: ListSKUs: %w", err)
	}

	result := make([]int64, len(responseObj.Skus))
	for i, sku := range responseObj.Skus {
		result[i] = int64(sku)
	}

	return result, nil
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (*domain.ProductSt, error) {
	responseObj, err := c.client.GetProduct(ctx, &product_service.GetProductRequest{
		Token: c.token,
		Sku:   uint32(sku),
	})
	if err != nil {
		return nil, fmt.Errorf("product_service: GetProduct: %w", err)
	}

	return &domain.ProductSt{
		Name:  responseObj.Name,
		Price: responseObj.Price,
	}, nil
}
