package models

type StockReserveSt struct {
	OrderId     int64
	WarehouseID int64
	Sku         uint32
	Count       uint64
}

type StockReserveListParsSt struct {
	OrderId     *int64
	WarehouseID *int64
	Sku         *uint32
}
