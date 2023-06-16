package models

type CartSt struct {
	Id         int64
	User       int64
	TotalPrice uint32
	Items      []*CartItemSt
}

type CartListParsSt struct {
}
