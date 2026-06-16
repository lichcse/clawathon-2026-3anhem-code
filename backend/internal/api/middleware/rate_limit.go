package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit returns a sliding-window rate limiter. If rdb is nil it falls back
// to allowing all requests (useful when Redis is unavailable).
func RateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		// Use user_id from JWT context if available; else fall back to IP.
		key := c.ClientIP()
		if uid, ok := c.Get("user_id"); ok {
			key = fmt.Sprintf("user:%v", uid)
		}
		redisKey := fmt.Sprintf("rate:%s", key)

		ctx := context.Background()
		pipe := rdb.Pipeline()
		incr := pipe.Incr(ctx, redisKey)
		pipe.Expire(ctx, redisKey, window)
		if _, err := pipe.Exec(ctx); err != nil {
			// Redis unavailable — fail open
			c.Next()
			return
		}

		count := incr.Val()
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxInt64(0, int64(limit)-count)))

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please slow down.",
				},
			})
			return
		}
		c.Next()
	}
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
