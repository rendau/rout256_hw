package handler

import (
	"context"

	"route256/checkout/internal/domain"
	"route256/checkout/pkg/proto/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	checkout_v1.UnimplementedCheckoutServer
	Model *domain.Domain
}

func New(model *domain.Domain) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) AddToCart(ctx context.Context, request *checkout_v1.AddToCartRequest) (*emptypb.Empty, error) {
	err := h.Model.AddToCart(ctx, request.User, request.Sku, uint16(request.Count))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteFromCart(ctx context.Context, request *checkout_v1.DeleteFromCartRequest) (*emptypb.Empty, error) {
	err := h.Model.DeleteFromCart(ctx, request.User, request.Sku, uint16(request.Count))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) ListCart(ctx context.Context, request *checkout_v1.ListCartRequest) (*checkout_v1.ListCartResponse, error) {
	cart, err := h.Model.ListCart(ctx, request.User)
	if err != nil {
		return nil, err
	}

	items := make([]*checkout_v1.ListCartResponseItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = &checkout_v1.ListCartResponseItem{
			Sku:   item.Sku,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		}
	}

	return &checkout_v1.ListCartResponse{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}, nil
}

func (h *Handler) Purchase(ctx context.Context, request *checkout_v1.PurchaseRequest) (*checkout_v1.PurchaseResponse, error) {
	orderID, err := h.Model.Purchase(ctx, request.User)
	if err != nil {
		return nil, err
	}

	return &checkout_v1.PurchaseResponse{
		OrderID: orderID,
	}, nil
}
