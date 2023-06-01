package purchase

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
	OrderID int64 `json:"orderID"`
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	orderId, err := h.Model.Purchase(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &Response{
		OrderID: orderId,
	}, nil
}
