package constant

const (
	OrderStatusNew             = "new"
	OrderStatusAwaitingPayment = "awaiting payment"
	OrderStatusFailed          = "failed"
	OrderStatusPayed           = "payed"
	OrderStatusCancelled       = "cancelled"
)

func OrderStatusIsValid(v string) bool {
	return v == OrderStatusNew ||
		v == OrderStatusAwaitingPayment ||
		v == OrderStatusFailed ||
		v == OrderStatusPayed ||
		v == OrderStatusCancelled
}
