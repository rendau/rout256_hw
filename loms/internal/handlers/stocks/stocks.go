package stocks

import (
	"context"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Request struct {
	SKU uint32 `json:"sku"`
}

type Response struct {
	Stocks []responseStock `json:"stocks"`
}

type responseStock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

func New(model *domain.Model) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	result, err := h.Model.Stocks(ctx, req.SKU)
	if err != nil {
		return nil, err
	}

	responseObj := &Response{
		Stocks: make([]responseStock, len(result.Stocks)),
	}

	for i, stock := range result.Stocks {
		responseObj.Stocks[i] = responseStock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		}
	}

	return responseObj, nil
}
