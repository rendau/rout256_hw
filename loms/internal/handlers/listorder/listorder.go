package listorder

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
	Status string         `json:"status"`
	User   int64          `json:"user"`
	Items  []responseItem `json:"items"`
}

type responseItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	result, err := h.Model.ListOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	responseObj := &Response{
		Status: result.Status,
		User:   result.User,
		Items:  make([]responseItem, len(result.Items)),
	}

	for i, item := range result.Items {
		responseObj.Items[i] = responseItem{
			SKU:   item.Sku,
			Count: item.Count,
		}
	}

	return responseObj, err
}
