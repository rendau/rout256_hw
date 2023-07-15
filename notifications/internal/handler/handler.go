package handler

import (
	"context"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/domain/models"
	"route256/notifications/pkg/proto/notifications_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	notifications_v1.UnimplementedNotificationsServer
	Model *domain.Domain
}

func New(model *domain.Domain) *Handler {
	return &Handler{Model: model}
}

func (h *Handler) ListOrderStatusEvent(ctx context.Context, request *notifications_v1.ListOrderStatusEventRequest) (*notifications_v1.ListOrderStatusEventResponse, error) {
	pars := &models.OrderStatusEventListParsSt{
		OrderID: request.OrderID,
	}
	if request.TsGTE != nil {
		t := request.TsGTE.AsTime()
		pars.TsGTE = &t
	}
	if request.TsLTE != nil {
		t := request.TsLTE.AsTime()
		pars.TsLTE = &t
	}
	result, err := h.Model.ListOrderStatusEvent(ctx, pars)
	if err != nil {
		return nil, err
	}

	items := make([]*notifications_v1.OrderStatusEvent, len(result))

	for i, item := range result {
		items[i] = &notifications_v1.OrderStatusEvent{
			Ts:      timestamppb.New(item.TS),
			OrderID: item.OrderID,
			Status:  item.Status,
		}
	}

	return &notifications_v1.ListOrderStatusEventResponse{
		Items: items,
	}, nil
}
