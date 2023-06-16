package pg

import (
	"context"
	"fmt"

	"route256/checkout/internal/domain/models"
	"route256/checkout/internal/repo/schema"
)

var (
	cartItemTableName  = "cart_item"
	cartItemAllColumns = []string{"cart_id", "sku", "cnt"}
)

func (r *St) CartItemList(ctx context.Context, pars models.CartItemListParsSt) ([]*models.CartItemSt, error) {
	query := r.sq.Select(cartItemAllColumns...).
		From(cartItemTableName)

	// filters
	if pars.CartId != nil {
		query = query.Where("cart_id = ?", *pars.CartId)
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

	result := make([]*models.CartItemSt, 0)

	for rows.Next() {
		item := &schema.CartItemSt{}

		err = rows.Scan(
			&item.CartId,
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

func (r *St) CartItemGet(ctx context.Context, cartId int64, sku uint32) (*models.CartItemSt, error) {
	query := r.sq.Select(cartItemAllColumns...).
		From(cartItemTableName).
		Where("cart_id = ? AND sku = ?", cartId, sku)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	item := &schema.CartItemSt{}

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(
		&item.CartId,
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

func (r *St) CartItemSet(ctx context.Context, obj *models.CartItemSt) error {
	rawSQL := `
		INSERT INTO cart_item (cart_id, sku, cnt)
		VALUES ($1, $2, $3)
		on conflict (cart_id, sku) do
		    update set cnt = $3
	`

	err := r.db.Exec(ctx, rawSQL, obj.CartId, obj.Sku, obj.Count)
	if err != nil {
		return fmt.Errorf("db.QueryRow: %w", err)
	}

	return nil
}

func (r *St) CartItemRemove(ctx context.Context, cartId int64, sku uint32) error {
	rawSQL := `
		DELETE FROM cart_item
		WHERE cart_id = $1 and sku = $2
	`

	err := r.db.Exec(ctx, rawSQL, cartId, sku)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}

func (r *St) CartItemRemoveAllForCartId(ctx context.Context, cartId int64) error {
	rawSQL := `
		DELETE FROM cart_item
		WHERE cart_id = $1
	`

	err := r.db.Exec(ctx, rawSQL, cartId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}
