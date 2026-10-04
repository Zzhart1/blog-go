package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimiter 固定窗口限流：同一 key 在 window 时间内最多 limit 次
type rateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	limit    int
	visitors map[string]*visit
}

type visit struct {
	count    int
	windowAt time.Time
}

func newRateLimiter(window time.Duration, limit int) *rateLimiter {
	return &rateLimiter{
		window:   window,
		limit:    limit,
		visitors: make(map[string]*visit),
	}
}

// allow 返回当前 key 是否允许继续请求
func (r *rateLimiter) allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	v, ok := r.visitors[key]
	if !ok || now.Sub(v.windowAt) >= r.window {
		// 新窗口开始
		r.visitors[key] = &visit{count: 1, windowAt: now}
		return true
	}
	v.count++
	return v.count <= r.limit
}

// LoginRateLimit 登录接口限流：每 IP 每分钟最多 10 次，防暴力破解
func LoginRateLimit() gin.HandlerFunc {
	rl := newRateLimiter(time.Minute, 10)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"msg": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}
