package pg

import (
	"context"
	"fmt"

	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/models"
	"route256/checkout/internal/repo/schema"
)

var (
	cartTableName  = "cart"
	cartAllColumns = []string{"id", "user_id"}
)

func (r *St) CartList(ctx context.Context, pars models.CartListParsSt) ([]*models.CartSt, error) {
	query := r.sq.Select(cartAllColumns...).
		From(cartTableName)

	// filters

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

	result := make([]*models.CartSt, 0)

	for rows.Next() {
		item := &schema.CartSt{}

		err = rows.Scan(
			&item.Id,
			&item.UserId,
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

func (r *St) CartGet(ctx context.Context, id int64) (*models.CartSt, error) {
	query := r.sq.Select(cartAllColumns...).
		From(cartTableName).
		Where("id = ?", id)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	item := &schema.CartSt{}

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(
		&item.Id,
		&item.UserId,
	)
	if err != nil {
		if r.db.ErrorIsNoRows(err) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("db.QueryRow: %w", err)
	}

	return item.ToModel(), nil
}

func (r *St) CartGetByUsrId(ctx context.Context, userId int64) (*models.CartSt, error) {
	query := r.sq.Select(cartAllColumns...).
		From(cartTableName).
		Where("user_id = ?", userId)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query.ToSql: %w", err)
	}

	item := &schema.CartSt{}

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(
		&item.Id,
		&item.UserId,
	)
	if err != nil {
		if r.db.ErrorIsNoRows(err) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("db.QueryRow: %w", err)
	}

	return item.ToModel(), nil
}

func (r *St) CartCreate(ctx context.Context, obj *models.CartSt) (int64, error) {
	rawSQL := `
		INSERT INTO cart (user_id)
		VALUES ($1)
		RETURNING id
	`

	var newId int64

	err := r.db.QueryRow(ctx, rawSQL, obj.User).Scan(&newId)
	if err != nil {
		return 0, fmt.Errorf("db.QueryRow: %w", err)
	}

	return newId, nil
}

func (r *St) CartRemove(ctx context.Context, id int64) error {
	rawSQL := `
		DELETE FROM cart
		WHERE id = $1
	`

	err := r.db.Exec(ctx, rawSQL, id)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	return nil
}
