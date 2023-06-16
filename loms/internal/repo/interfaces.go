package repo

import (
	"context"

	"route256/loms/internal/domain/models"
)

type Repo interface {
	// Stock
	StockList(ctx context.Context, pars models.StockListParsSt, lock bool) ([]*models.StockSt, error)
	StockPut(ctx context.Context, warehouseID int64, sku uint32, count uint64) error
	StockPull(ctx context.Context, warehouseID int64, sku uint32, count uint64) error
	StockGet(ctx context.Context, warehouseID int64, sku uint32) (*models.StockSt, error)
	StockRemove(ctx context.Context, warehouseID int64, sku uint32) error

	// StockReserve
	StockReserveList(ctx context.Context, pars models.StockReserveListParsSt) ([]*models.StockReserveSt, error)
	StockReservePut(ctx context.Context, orderID, warehouseID int64, sku uint32, count uint64) error
	StockReservePull(ctx context.Context, orderID, warehouseID int64, sku uint32, count uint64) error
	StockReserveGet(ctx context.Context, orderID, warehouseID int64, sku uint32) (*models.StockReserveSt, error)
	StockReserveRemove(ctx context.Context, orderID, warehouseID int64, sku uint32) error

	// Order
	OrderList(ctx context.Context, pars models.OrderListParsSt) ([]*models.OrderListSt, error)
	OrderGet(ctx context.Context, id int64) (*models.OrderSt, error)
	OrderCreate(ctx context.Context, obj *models.OrderSt) (int64, error)
	OrderRemove(ctx context.Context, id int64) error
	OrderGetItems(ctx context.Context, id int64) ([]*models.OrderItemSt, error)
	OrderSetStatus(ctx context.Context, id int64, v string) error
}
