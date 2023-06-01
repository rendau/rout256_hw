package loms

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksResponse struct {
	Stocks []struct {
		WarehouseID int64  `json:"warehouseID"`
		Count       uint64 `json:"count"`
	} `json:"stocks"`
}

type CreateOrderRequest struct {
	User  int64                    `json:"user"`
	Items []CreateOrderRequestItem `json:"items"`
}

type CreateOrderRequestItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type CreateOrderResponse struct {
	OrderID int64 `json:"orderID"`
}
