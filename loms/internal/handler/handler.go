package handler

import (
	"context"

	"route256/loms/internal/domain"
	"route256/loms/internal/domain/models"
	"route256/loms/pkg/proto/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	loms_v1.UnimplementedLomsServer
	Model *domain.Domain
}

func New(model *domain.Domain) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) CancelOrder(ctx context.Context, request *loms_v1.CancelOrderRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.Model.CancelOrder(ctx, request.OrderID)
}

func (h *Handler) CreateOrder(ctx context.Context, request *loms_v1.CreateOrderRequest) (*loms_v1.CreateOrderResponse, error) {
	items := make([]*models.OrderItemSt, len(request.Items))
	for i, item := range request.Items {
		items[i] = &models.OrderItemSt{
			Sku:   item.Sku,
			Count: uint16(item.Count),
		}
	}

	orderID, err := h.Model.CreateOrder(ctx, request.User, items)

	return &loms_v1.CreateOrderResponse{OrderID: orderID}, err
}

func (h *Handler) ListOrder(ctx context.Context, request *loms_v1.ListOrderRequest) (*loms_v1.ListOrderResponse, error) {
	result, err := h.Model.ListOrder(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}

	responseObj := &loms_v1.ListOrderResponse{
		Status: result.Status,
		User:   result.User,
		Items:  make([]*loms_v1.Order, len(result.Items)),
	}

	for i, item := range result.Items {
		responseObj.Items[i] = &loms_v1.Order{
			Sku:   item.Sku,
			Count: uint32(item.Count),
		}
	}

	return responseObj, err
}

func (h *Handler) OrderPayed(ctx context.Context, request *loms_v1.OrderPayedRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.Model.OrderPayed(ctx, request.OrderID)
}

func (h *Handler) Stocks(ctx context.Context, request *loms_v1.StocksRequest) (*loms_v1.StocksResponse, error) {
	items, err := h.Model.Stocks(ctx, request.Sku)
	if err != nil {
		return nil, err
	}

	result := make([]*loms_v1.Stock, len(items))

	for i, item := range items {
		result[i] = &loms_v1.Stock{
			WarehouseID: item.WarehouseID,
			Count:       item.Count,
		}
	}

	return &loms_v1.StocksResponse{Stocks: result}, nil
}

func (h *Handler) StockAdd(ctx context.Context, request *loms_v1.StockAddRequest) (*emptypb.Empty, error) {
	err := h.Model.StockAdd(ctx, request.WarehouseId, request.Sku, request.Count)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) StockRemove(ctx context.Context, request *loms_v1.StockRemoveRequest) (*emptypb.Empty, error) {
	err := h.Model.StockRemove(ctx, request.WarehouseId, request.Sku)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
