package schema

import (
	"route256/loms/internal/domain/models"
)

// stock

type StockSt struct {
	WarehouseID int64  `db:"warehouse_id"`
	Sku         uint32 `db:"sku"`
	Count       uint64 `db:"cnt"`
}

func (o *StockSt) ToModel() *models.StockSt {
	return &models.StockSt{
		WarehouseID: o.WarehouseID,
		Sku:         o.Sku,
		Count:       o.Count,
	}
}
