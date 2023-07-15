package pg

import (
	"context"
	"fmt"
	"route256/notifications/internal/domain/models"
	"route256/notifications/internal/repo/schema"
)

var (
	orderStatusEventTableName  = "order_status_event"
	orderStatusEventAllColumns = []string{"ts", "order_id", "status"}
)

func (r *St) OrderStatusEventCreate(ctx context.Context, obj *models.OrderStatusEventSt) error {
	rawSQL := `
		INSERT INTO ` + orderStatusEventTableName + ` (ts, order_id, status)
		VALUES ($1, $2, $3)
	`

	err := r.db.Exec(ctx, rawSQL, obj.TS, obj.OrderID, obj.Status)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) OrderStatusEventList(ctx context.Context, pars *models.OrderStatusEventListParsSt) ([]*models.OrderStatusEventSt, error) {
	query := r.sq.Select(orderStatusEventAllColumns...).
		From(orderStatusEventTableName)

	// filters
	if pars.TsGTE != nil {
		query = query.Where("ts >= ?", *pars.TsGTE)
	}
	if pars.TsLTE != nil {
		query = query.Where("ts <= ?", *pars.TsLTE)
	}
	if pars.OrderID != nil {
		query = query.Where("order_id = ?", *pars.OrderID)
	}

	query.OrderBy("ts")

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
	result := make([]*models.OrderStatusEventSt, 0)
	for rows.Next() {
		item := &schema.OrderStatusEventSt{}

		err = rows.Scan(
			&item.TS,
			&item.OrderId,
			&item.Status,
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
