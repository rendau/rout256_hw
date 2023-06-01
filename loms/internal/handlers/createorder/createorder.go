package createorder

import (
	"context"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Request struct {
	User  int64 `json:"user"`
	Items []struct {
		SKU   uint32 `json:"sku"`
		Count uint16 `json:"count"`
	} `json:"items"`
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	reqObj := &domain.CreateOrderRequestSt{
		User:  req.User,
		Items: make([]domain.CreateOrderRequestItemSt, len(req.Items)),
	}
	for i, item := range req.Items {
		reqObj.Items[i] = domain.CreateOrderRequestItemSt{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}

	orderID, err := h.Model.CreateOrder(ctx, reqObj)

	return &Response{OrderID: orderID}, err
}
