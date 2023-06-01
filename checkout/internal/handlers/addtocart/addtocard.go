package addtocart

import (
	"context"
	"route256/checkout/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct {
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	if req == nil {
		req = &Request{}
	}

	err := h.Model.AddToCart(ctx, req.User, req.SKU, req.Count)
	if err != nil {
		return nil, err
	}

	return &Response{}, nil
}
