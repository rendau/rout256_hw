package schema

import (
	"route256/checkout/internal/domain/models"
)

type CartItemSt struct {
	CartId int64  `db:"cart_id"`
	Sku    uint32 `db:"sku"`
	Count  uint16 `db:"count"`
}

type CartItemListParsSt struct {
	CartId *int64
}

func (s *CartItemSt) ToModel() *models.CartItemSt {
	return &models.CartItemSt{
		CartId: s.CartId,
		Sku:    s.Sku,
		Count:  s.Count,
	}
}
