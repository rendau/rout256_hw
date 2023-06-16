package pg

import (
	"context"
	"fmt"

	"route256/loms/internal/domain/models"
	"route256/loms/internal/repo/schema"
)

var (
	stockReserveTableName  = "stock_reserve"
	stockReserveAllColumns = []string{"order_id", "warehouse_id", "sku", "cnt"}
)

func (r *St) StockReserveList(ctx context.Context, pars models.StockReserveListParsSt) ([]*models.StockReserveSt, error) {
	query := r.sq.Select(stockReserveAllColumns...).
		From(stockReserveTableName)

	// filters

	if pars.OrderId != nil {
		query = query.Where("order_id = ?", *pars.OrderId)
	}
	if pars.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *pars.WarehouseID)
	}
	if pars.Sku != nil {
		query = query.Where("sku = ?", *pars.Sku)
	}

	// raw sql

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	// send query

	rows, err := r.db.Query(ctx, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	// scan

	result := make([]*models.StockReserveSt, 0)

	for rows.Next() {
		item := &schema.StockReserveSt{}

		err = rows.Scan(
			&item.OrderId,
			&item.WarehouseID,
			&item.Sku,
			&item.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		result = append(result, item.ToModel())
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return result, nil
}

func (r *St) StockReservePut(ctx context.Context, orderID, warehouseID int64, sku uint32, count uint64) error {
	rawSQL := `
		INSERT INTO stock_reserve (order_id, warehouse_id, sku, cnt)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (order_id, warehouse_id, sku) DO UPDATE
			SET cnt = stock_reserve.cnt + $4
	`

	err := r.db.Exec(ctx, rawSQL, orderID, warehouseID, sku, count)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) StockReservePull(ctx context.Context, orderID, warehouseID int64, sku uint32, count uint64) error {
	rawSQL := `
		UPDATE stock_reserve
		SET cnt = stock_reserve.cnt - $4
		WHERE order_id = $1 AND warehouse_id = $2 AND sku = $3
	`

	err := r.db.Exec(ctx, rawSQL, orderID, warehouseID, sku, count)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) StockReserveGet(ctx context.Context, orderID, warehouseID int64, sku uint32) (*models.StockReserveSt, error) {
	query := r.sq.Select(stockReserveAllColumns...).
		From(stockReserveTableName).
		Where("order_id = ? AND warehouse_id = ? AND sku = ?", orderID, warehouseID, sku)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	item := &schema.StockReserveSt{}

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(
		&item.OrderId,
		&item.WarehouseID,
		&item.Sku,
		&item.Count,
	)
	if err != nil {
		if r.db.ErrorIsNoRows(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRow: %w", err)
	}

	return item.ToModel(), nil
}

func (r *St) StockReserveRemove(ctx context.Context, orderID, warehouseID int64, sku uint32) error {
	rawSQL := `
		DELETE FROM stock_reserve
		WHERE order_id = $1 AND warehouse_id = $2 AND sku = $3
	`

	err := r.db.Exec(ctx, rawSQL, orderID, warehouseID, sku)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}
