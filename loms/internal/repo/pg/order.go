package pg

import (
	"context"
	"fmt"

	"route256/loms/internal/domain/models"
	"route256/loms/internal/repo/schema"
)

var (
	orderTableName  = "ord"
	orderAllColumns = []string{"id", "user_id", "status"}
)

func (r *St) OrderList(ctx context.Context, pars models.OrderListParsSt) ([]*models.OrderListSt, error) {
	query := r.sq.Select(orderAllColumns...).
		From(orderTableName)

	// filters

	if pars.User != nil {
		query = query.Where("user_id = ?", *pars.User)
	}
	if pars.Status != nil {
		query = query.Where("status = ?", *pars.Status)
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

	result := make([]*models.OrderListSt, 0)

	for rows.Next() {
		item := &schema.OrderListSt{}

		err = rows.Scan(
			&item.ID,
			&item.User,
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

func (r *St) OrderGet(ctx context.Context, id int64) (*models.OrderSt, error) {
	query := r.sq.Select(orderAllColumns...).
		From(orderTableName).
		Where("id = ?", id)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	item := &schema.OrderSt{}

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(
		&item.ID,
		&item.User,
		&item.Status,
	)
	if err != nil {
		if r.db.ErrorIsNoRows(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRow: %w", err)
	}

	return item.ToModel(), nil
}

func (r *St) OrderCreate(ctx context.Context, obj *models.OrderSt) (int64, error) {
	rawSQL := `
		INSERT INTO ord (user_id, status)
		VALUES ($1, $2)
		RETURNING id
	`

	var newId int64

	err := r.db.QueryRow(ctx, rawSQL, obj.User, obj.Status).Scan(&newId)
	if err != nil {
		return 0, fmt.Errorf("db.QueryRow: %w", err)
	}

	return newId, nil
}

func (r *St) OrderRemove(ctx context.Context, id int64) error {
	rawSQL := `
		DELETE FROM ord
		WHERE id = $1
	`

	err := r.db.Exec(ctx, rawSQL, id)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) OrderGetItems(ctx context.Context, id int64) ([]*models.OrderItemSt, error) {
	query := r.sq.Select("sku", "sum(cnt)").
		From(stockReserveTableName).
		Where("order_id = ?", id).
		GroupBy("sku")

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

	result := make([]*models.OrderItemSt, 0)

	for rows.Next() {
		item := &schema.OrderItemSt{}

		err = rows.Scan(
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

func (r *St) OrderSetStatus(ctx context.Context, id int64, v string) error {
	rawSQL := `
		UPDATE ord
		SET status = $1
		WHERE id = $2
	`

	err := r.db.Exec(ctx, rawSQL, v, id)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}
