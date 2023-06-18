package ratelimiter

import (
	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	limiter ratelimit.Limiter
}

func NewRateLimiter(n int) *RateLimiter {
	return &RateLimiter{
		limiter: ratelimit.New(n),
	}
}

func (r *RateLimiter) Take() {
	r.limiter.Take()
}
