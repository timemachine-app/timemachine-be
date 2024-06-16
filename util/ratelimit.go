package util

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/timemachine-app/timemachine-be/internal/config"
)

// TODO: Use Redis when you start horizontal scale
var rateLimitStore = struct {
	sync.Mutex
	clients map[string][]int64
}{
	clients: make(map[string][]int64),
}

// Rate limiting middleware
func RateLimitMiddleware(ratelimitConfig config.RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		currentTime := time.Now().Unix()

		rateLimitStore.Lock()
		defer rateLimitStore.Unlock()

		// Initialize the request times slice if not present
		if _, exists := rateLimitStore.clients[clientIP]; !exists {
			rateLimitStore.clients[clientIP] = []int64{}
		}

		// Append the current request time
		rateLimitStore.clients[clientIP] = append(rateLimitStore.clients[clientIP], currentTime)

		// Remove timestamps older than the time window
		validTime := currentTime - ratelimitConfig.WindowInSec
		validRequests := []int64{}
		for _, t := range rateLimitStore.clients[clientIP] {
			if t > validTime {
				validRequests = append(validRequests, t)
			}
		}
		rateLimitStore.clients[clientIP] = validRequests

		// Check if the request count exceeds the limit
		if len(validRequests) > ratelimitConfig.RateLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
