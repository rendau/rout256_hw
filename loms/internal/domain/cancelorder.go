package domain

import (
	"context"
	"fmt"

	"route256/libs/constant"
	"route256/loms/internal/domain/models"
)

func (d *Domain) CancelOrderValidate(orderID int64) error {
	if orderID <= 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (d *Domain) CancelOrder(ctx context.Context, orderID int64) error {
	var err error

	// validate
	if err = d.CancelOrderValidate(orderID); err != nil {
		return err
	}

	// get order
	order, err := d.repo.OrderGet(ctx, orderID)
	if err != nil {
		return fmt.Errorf("OrderGet: %w", err)
	}

	if order.Status == constant.OrderStatusCancelled {
		return nil
	}

	err = d.db.TransactionFn(ctx, func(ctx context.Context) error {
		// get reserve items
		reserveItems, err := d.repo.StockReserveList(ctx, models.StockReserveListParsSt{
			OrderId: &orderID,
		})
		if err != nil {
			return fmt.Errorf("StockReserveList: %w", err)
		}

		// release reserve items
		for _, item := range reserveItems {
			// put stock
			err = d.repo.StockPut(ctx, item.WarehouseID, item.Sku, item.Count)
			if err != nil {
				return fmt.Errorf("StockPut: %w", err)
			}

			// remove reserve
			err = d.repo.StockReserveRemove(ctx, orderID, item.WarehouseID, item.Sku)
			if err != nil {
				return fmt.Errorf("StockReserveRemove: %w", err)
			}
		}

		// update order status
		err = d.repo.OrderSetStatus(ctx, orderID, constant.OrderStatusCancelled)
		if err != nil {
			return fmt.Errorf("OrderSetStatus: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// send notification
	err = d.NotificationSendOrderStatusChange(models.NotificationOrderStatusChangeSt{
		OrderID: orderID,
		Status:  constant.OrderStatusCancelled,
	})
	if err != nil {
		return err
	}

	return nil
}
