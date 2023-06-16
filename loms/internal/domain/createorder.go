package domain

import (
	"context"
	"math"

	"route256/libs/constant"
	"route256/loms/internal/domain/models"
)

func (d *Domain) CreateOrderValidate(userId int64, items []*models.OrderItemSt) error {
	if userId <= 0 {
		return ErrUserNotFound
	}

	// validate items
	for _, item := range items {
		if item.Sku == 0 {
			return ErrSkuRequired
		}
		if item.Count == 0 {
			return ErrCountRequired
		}
	}

	return nil
}

func (d *Domain) CreateOrder(ctx context.Context, userId int64, items []*models.OrderItemSt) (int64, error) {
	// validate
	if err := d.CreateOrderValidate(userId, items); err != nil {
		return 0, err
	}

	var stockFilterCntGT uint64 = 0
	var newId int64

	err := d.db.TransactionFn(ctx, func(ctx context.Context) error {
		// get stock items
		stockItems, err := d.repo.StockList(ctx, models.StockListParsSt{
			CountGT: &stockFilterCntGT,
		}, true)
		if err != nil {
			return err
		}

		// distribute items in stock
		reserveItems, err := d.distributeItemsInStock(items, stockItems)
		if err != nil {
			return err
		}

		// create order
		newId, err = d.repo.OrderCreate(ctx, &models.OrderSt{
			User:   userId,
			Status: constant.OrderStatusNew,
		})
		if err != nil {
			return err
		}

		// reserve
		for _, item := range reserveItems {
			// pull stock
			err = d.repo.StockPull(ctx, item.WarehouseID, item.Sku, item.Count)
			if err != nil {
				return err
			}

			// put reserve
			err = d.repo.StockReservePut(ctx, newId, item.WarehouseID, item.Sku, item.Count)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return newId, err
}

func (d *Domain) distributeItemsInStock(items []*models.OrderItemSt, stock []*models.StockSt) ([]*models.StockReserveSt, error) {
	var result []*models.StockReserveSt

	for _, item := range items {
		count := item.Count

		for _, stockItem := range stock {
			if stockItem.Sku != item.Sku {
				continue
			}
			if stockItem.Count <= 0 {
				continue
			}

			stockReserve := math.Min(float64(count), float64(stockItem.Count))

			if stockReserve > 0 {
				result = append(result, &models.StockReserveSt{
					WarehouseID: stockItem.WarehouseID,
					Sku:         stockItem.Sku,
					Count:       uint64(stockReserve),
				})

				count -= uint16(stockReserve)
			}

			if count <= 0 {
				break
			}
		}

		if count > 0 {
			return nil, ErrStockInsufficient
		}
	}

	return result, nil
}
