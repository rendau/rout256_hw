package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type St struct {
	r *redis.Client
}

func New(url, psw string, db int) *St {
	return &St{
		r: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: psw,
			DB:       db,
		}),
	}
}

func (c *St) Get(ctx context.Context, key string) ([]byte, bool, error) {
	data, err := c.r.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis: fail to 'get': %w", err)
	}

	return data, true, nil
}

func (c *St) GetJsonObj(ctx context.Context, key string, dst any) (bool, error) {
	dataRaw, ok, err := c.Get(ctx, key)
	if err != nil || !ok {
		return ok, err
	}

	err = json.Unmarshal(dataRaw, dst)
	if err != nil {
		return false, fmt.Errorf("redis: fail to unmarshal json: %w", err)
	}

	return true, nil
}

func (c *St) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	err := c.r.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis: fail to 'set': %w", err)
	}

	return nil
}

func (c *St) SetJsonObj(ctx context.Context, key string, value any, expiration time.Duration) error {
	dataRaw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, dataRaw, expiration)
}

func (c *St) Del(ctx context.Context, key string) error {
	err := c.r.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis: fail to 'del': %w", err)
	}

	return nil
}
