package domain

import (
	"context"
	"time"
)

//go:generate mockery --name IMessanger --output ./mocks --filename messanger.go
type IMessanger interface {
	Send(msg string) error
}

//go:generate mockery --name ICache --output ./mocks --filename cache.go
type ICache interface {
	GetJsonObj(ctx context.Context, key string, dst any) (bool, error)
	SetJsonObj(ctx context.Context, key string, value any, expiration time.Duration) error
}
