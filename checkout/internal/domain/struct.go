package domain

type StockSt struct {
	WarehouseID int64
	Count       uint64
}

type CartSt struct {
	Items      []*CartItemSt
	TotalPrice uint32
}

type CartItemSt struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type ProductSt struct {
	Name  string
	Price uint32
}
