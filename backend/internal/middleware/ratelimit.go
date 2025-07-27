package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimitMiddleware limits the number of requests per client
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		// Clean up old requests
		now := time.Now()
		cutoff := now.Add(-window)
		
		var validRequests []time.Time
		for _, t := range limiter.requests[clientIP] {
			if t.After(cutoff) {
				validRequests = append(validRequests, t)
			}
		}
		
		limiter.requests[clientIP] = validRequests

		// Check if limit is exceeded
		if len(validRequests) >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		// Add current request
		limiter.requests[clientIP] = append(limiter.requests[clientIP], now)
		c.Next()
	}
}