package schema

import (
	"route256/loms/internal/domain/models"
)

type StockReserveSt struct {
	OrderId     int64  `db:"order_id"`
	WarehouseID int64  `db:"warehouse_id"`
	Sku         uint32 `db:"sku"`
	Count       uint64 `db:"cnt"`
}

func (o *StockReserveSt) ToModel() *models.StockReserveSt {
	return &models.StockReserveSt{
		OrderId:     o.OrderId,
		WarehouseID: o.WarehouseID,
		Sku:         o.Sku,
		Count:       o.Count,
	}
}
