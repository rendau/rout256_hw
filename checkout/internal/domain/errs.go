package domain

type Err string

func (e Err) Error() string {
	return string(e)
}

// error constants

const (
	ErrUserNotFound      = Err("user_not_found")
	ErrStockInsufficient = Err("stock_insufficient")
	ErrSkuRequired       = Err("sku_required")
	ErrCountRequired     = Err("count_required")
	ErrCartNotFound      = Err("cart_not_found")
)
