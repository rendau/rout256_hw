package pg

import (
	"context"
	"fmt"

	"route256/loms/internal/domain/models"
	"route256/loms/internal/repo/schema"
)

var (
	stockTableName  = "stock"
	stockAllColumns = []string{"warehouse_id", "sku", "cnt"}
)

func (r *St) StockList(ctx context.Context, pars models.StockListParsSt, lock bool) ([]*models.StockSt, error) {
	query := r.sq.Select(stockAllColumns...).
		From(stockTableName)

	if lock {
		query = query.Suffix("FOR UPDATE")
	}

	// filters

	if pars.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *pars.WarehouseID)
	}
	if pars.Sku != nil {
		query = query.Where("sku = ?", *pars.Sku)
	}
	if pars.CountGT != nil {
		query = query.Where("cnt > ?", *pars.CountGT)
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

	result := make([]*models.StockSt, 0)

	for rows.Next() {
		item := &schema.StockSt{}

		err = rows.Scan(
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

func (r *St) StockPut(ctx context.Context, warehouseID int64, sku uint32, count uint64) error {
	rawSQL := `
		INSERT INTO stock (warehouse_id, sku, cnt)
		VALUES ($1, $2, $3)
		ON CONFLICT (warehouse_id, sku) DO UPDATE
			SET cnt = stock.cnt + $3
	`

	err := r.db.Exec(ctx, rawSQL, warehouseID, sku, count)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) StockPull(ctx context.Context, warehouseID int64, sku uint32, count uint64) error {
	rawSQL := `
		UPDATE stock
		SET cnt = stock.cnt - $3
		WHERE warehouse_id = $1 AND sku = $2
	`

	err := r.db.Exec(ctx, rawSQL, warehouseID, sku, count)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) StockGet(ctx context.Context, warehouseID int64, sku uint32) (*models.StockSt, error) {
	query := r.sq.Select(stockAllColumns...).
		From(stockTableName).
		Where("warehouse_id = ? AND sku = ?", warehouseID, sku)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	item := &schema.StockSt{}

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(
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

func (r *St) StockRemove(ctx context.Context, warehouseID int64, sku uint32) error {
	rawSQL := `
		DELETE FROM stock
		WHERE warehouse_id = $1 AND sku = $2
	`

	err := r.db.Exec(ctx, rawSQL, warehouseID, sku)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}
