package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	GetJsonObj(ctx context.Context, key string, dst any) (bool, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	SetJsonObj(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}
