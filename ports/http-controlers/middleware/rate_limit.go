package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // tokens per second
	burst    int           // max tokens
	cleanup  time.Duration // cleanup interval for old entries
}

type visitor struct {
	tokens     float64
	lastUpdate time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  time.Minute * 5,
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// cleanupLoop removes stale entries periodically
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastUpdate) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:     float64(rl.burst - 1),
			lastUpdate: now,
		}
		return true
	}

	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(v.lastUpdate).Seconds()
	v.tokens += elapsed * float64(rl.rate)
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}
	v.lastUpdate = now

	if v.tokens >= 1 {
		v.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware returns a Gin middleware for rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// Get rate limit from environment, default to 100 requests per second
	rateStr := os.Getenv("RATE_LIMIT_RPS")
	rate := 100
	if rateStr != "" {
		if r, err := strconv.Atoi(rateStr); err == nil && r > 0 {
			rate = r
		}
	}

	// Get burst from environment, default to 200
	burstStr := os.Getenv("RATE_LIMIT_BURST")
	burst := 200
	if burstStr != "" {
		if b, err := strconv.Atoi(burstStr); err == nil && b > 0 {
			burst = b
		}
	}

	limiter := NewRateLimiter(rate, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
