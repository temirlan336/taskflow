package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redis  *redis.Client
	window time.Duration
	limit  int64
}

func NewRateLimiter(r *redis.Client, window time.Duration, limit int64) *RateLimiter {
	return &RateLimiter{
		redis:  r,
		window: window,
		limit:  limit,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	count, err := r.redis.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		if err := r.redis.Expire(ctx, key, r.window).Err(); err != nil {
			return false, err
		}
	}
	return count <= r.limit, nil

}

func RateLimitMiddleware(limit *RateLimiter, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey, ok := GetAPIKeyFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			redisKey := "rate:key" + apiKey + action

			allowed, err := limit.Allow(r.Context(), redisKey)
			if err != nil {
				http.Error(w, "rate limiter unavailable", http.StatusServiceUnavailable)
				return
			}
			if !allowed {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
