package repo

import (
	"context"
	"route256/notifications/internal/domain/models"
)

type Repo interface {
	// OrderStatusEvent
	OrderStatusEventCreate(ctx context.Context, obj *models.OrderStatusEventSt) error
	OrderStatusEventList(ctx context.Context, pars *models.OrderStatusEventListParsSt) ([]*models.OrderStatusEventSt, error)
}
