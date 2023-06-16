package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/pressly/goose/v3"

	"route256/libs/db"
)

type St struct {
	dsn string
	Con *pgxpool.Pool
}

func New(dsn string) (*St, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	cfg.MaxConns = 100
	cfg.MinConns = 2
	cfg.LazyConnect = true

	dbPool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ConnectConfig: %w", err)
	}

	return &St{
		dsn: dsn,
		Con: dbPool,
	}, nil
}

// transaction

func (d *St) RenewTransaction(ctx context.Context) (context.Context, error) {
	var err error

	err = d.commitContextTransaction(ctx)
	if err != nil {
		return ctx, err
	}

	return d.contextWithTransaction(ctx)
}

func (d *St) TransactionFn(ctx context.Context, f func(context.Context) error) error {
	var err error

	if ctx == nil {
		ctx = context.Background()
	}

	if ctx, err = d.contextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { d.rollbackContextTransaction(ctx) }()

	err = f(ctx)
	if err != nil {
		return err
	}

	return d.commitContextTransaction(ctx)
}

func (d *St) getContextTransaction(ctx context.Context) pgx.Tx {
	if v := ctx.Value(transactionCtxKey); v != nil {
		tr, ok := v.(pgx.Tx)
		if ok {
			return tr
		}
	}
	return nil
}

func (d *St) contextWithTransaction(ctx context.Context) (context.Context, error) {
	tx, err := d.Con.Begin(ctx)
	if err != nil {
		return ctx, fmt.Errorf("tx.Begin: %w", err)
	}

	return context.WithValue(ctx, transactionCtxKey, tx), nil
}

func (d *St) commitContextTransaction(ctx context.Context) error {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return nil
	}

	err := tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

func (d *St) rollbackContextTransaction(ctx context.Context) {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return
	}

	_ = tx.Rollback(ctx)
}

// query

func (d *St) Exec(ctx context.Context, sql string, args ...any) error {
	if tx := d.getContextTransaction(ctx); tx != nil {
		_, err := tx.Exec(ctx, sql, args...)
		return err
	}

	_, err := d.Con.Exec(ctx, sql, args...)
	return err
}

func (d *St) Query(ctx context.Context, sql string, args ...any) (db.Rows, error) {
	var err error
	var rows pgx.Rows

	if tx := d.getContextTransaction(ctx); tx != nil {
		rows, err = tx.Query(ctx, sql, args...)
	} else {
		rows, err = d.Con.Query(ctx, sql, args...)
	}

	return &rowsSt{Rows: rows}, err
}

func (d *St) QueryRow(ctx context.Context, sql string, args ...any) db.Row {
	var row pgx.Row

	if tx := d.getContextTransaction(ctx); tx != nil {
		row = tx.QueryRow(ctx, sql, args...)
	} else {
		row = d.Con.QueryRow(ctx, sql, args...)
	}

	return &rowSt{Row: row}
}

// handle error

func (d *St) ErrorIsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows)
}

func (d *St) Migrate(dir string) error {
	con, err := sql.Open("pgx", d.dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer con.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err = goose.Up(con, dir); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}
