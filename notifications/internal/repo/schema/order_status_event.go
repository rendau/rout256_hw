package schema

import (
	"route256/notifications/internal/domain/models"
	"time"
)

type OrderStatusEventSt struct {
	TS      time.Time `db:"ts"`
	OrderId int64     `db:"order_id"`
	Status  string    `db:"status"`
}

func (o *OrderStatusEventSt) ToModel() *models.OrderStatusEventSt {
	return &models.OrderStatusEventSt{
		TS:      o.TS,
		OrderID: o.OrderId,
		Status:  o.Status,
	}
}
