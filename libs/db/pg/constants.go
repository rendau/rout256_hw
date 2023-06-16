package pg

type transactionCtxKeyType bool

const (
	transactionCtxKey = transactionCtxKeyType(true)
)
