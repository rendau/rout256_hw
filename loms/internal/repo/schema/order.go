package schema

import (
	"route256/loms/internal/domain/models"
)

type OrderSt struct {
	ID     int64  `db:"id"`
	User   int64  `db:"user_id"`
	Status string `db:"status"`
}

type OrderItemSt struct {
	Sku   uint32 `db:"sku"`
	Count uint16 `db:"cnt"`
}

type OrderListSt struct {
	ID     int64  `db:"id"`
	User   int64  `db:"user_id"`
	Status string `db:"status"`
}

func (o *OrderSt) ToModel() *models.OrderSt {
	return &models.OrderSt{
		ID:     o.ID,
		User:   o.User,
		Status: o.Status,
	}
}

func (o *OrderListSt) ToModel() *models.OrderListSt {
	return &models.OrderListSt{
		ID:     o.ID,
		User:   o.User,
		Status: o.Status,
	}
}

func (o *OrderItemSt) ToModel() *models.OrderItemSt {
	return &models.OrderItemSt{
		Sku:   o.Sku,
		Count: o.Count,
	}
}
