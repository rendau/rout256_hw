package orderpayed

import (
	"context"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Request struct {
	OrderID int64 `json:"orderID"`
}

type Response struct {
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	err := h.Model.OrderPayed(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	return &Response{}, err
}
