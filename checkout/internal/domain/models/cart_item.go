package models

type CartItemSt struct {
	CartId int64
	Sku    uint32
	Count  uint16
	Name   string
	Price  uint32
}

type CartItemListParsSt struct {
	CartId *int64
}
