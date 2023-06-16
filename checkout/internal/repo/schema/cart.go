package schema

import (
	"route256/checkout/internal/domain/models"
)

type CartSt struct {
	Id     int64 `db:"id"`
	UserId int64 `db:"user_id"`
}

type CartListParsSt struct {
}

func (s *CartSt) ToModel() *models.CartSt {
	return &models.CartSt{
		Id:   s.Id,
		User: s.UserId,
	}
}
