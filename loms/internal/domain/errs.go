package domain

type Err string

func (e Err) Error() string {
	return string(e)
}

// error constants

const (
	ErrUserNotFound      = Err("user_not_found")
	ErrOrderNotFound     = Err("order_not_found")
	ErrSkuRequired       = Err("sku_required")
	ErrCountRequired     = Err("count_required")
	ErrStockInsufficient = Err("stock_insufficient")
)
