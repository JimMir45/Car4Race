package middleware

import (
	"net/http"
	"sync"
	"time"

	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

// 简单的内存限流器
type rateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

var limiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100,             // 每个IP每分钟最多100次请求
	window:   time.Minute,
}

// RateLimit 频率限制中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		windowStart := now.Add(-limiter.window)

		// 清理过期记录
		var valid []time.Time
		for _, t := range limiter.requests[ip] {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		limiter.requests[ip] = valid

		// 检查是否超过限制
		if len(valid) >= limiter.limit {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 记录本次请求
		limiter.requests[ip] = append(limiter.requests[ip], now)

		c.Next()
	}
}
