package listcart

import (
	"context"
	"route256/checkout/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Request struct {
	User int64 `json:"user"`
}

type Response struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"totalPrice"`
}

type CartItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	result, err := h.Model.ListCart(ctx, req.User)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Items:      make([]CartItem, len(result.Items)),
		TotalPrice: result.TotalPrice,
	}

	for i, item := range result.Items {
		resp.Items[i] = CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		}
	}

	return resp, nil
}
