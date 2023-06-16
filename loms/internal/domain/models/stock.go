package models

type StockSt struct {
	WarehouseID int64
	Sku         uint32
	Count       uint64
}

type StockListParsSt struct {
	WarehouseID *int64
	Sku         *uint32
	CountGT     *uint64
}
