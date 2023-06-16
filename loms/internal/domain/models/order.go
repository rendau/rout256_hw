package models

type OrderSt struct {
	ID     int64
	User   int64
	Status string
	Items  []*OrderItemSt
}

type OrderItemSt struct {
	Sku   uint32
	Count uint16
}

type OrderListSt struct {
	ID     int64
	User   int64
	Status string
}

type OrderListParsSt struct {
	User   *int64
	Status *string
}
