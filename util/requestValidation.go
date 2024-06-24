package util

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/timemachine-app/timemachine-be/internal/config"
	"github.com/timemachine-app/timemachine-be/superbase"
)

// TODO: Use Redis when you start horizontal scale
var rateLimitStore = struct {
	sync.Mutex
	clients map[string][]int64
}{
	clients: make(map[string][]int64),
}

// Function to validate JWT token and return userId or nil
func validateToken(tokenString, jwtSecret string) *string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used to sign the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["sub"].(string); ok {
			return &userId
		}
	}

	return nil
}

// Rate limiting + token validation middleware
func ValidationMiddleware(
	ratelimitConfig config.RateLimitConfig, jwtSecret string, superbaseClient *superbase.SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for /health endpoint
		if strings.HasPrefix(c.Request.URL.Path, "/health") {
			c.Next()
			return
		}

		clientIdentifier := c.ClientIP() // Default to IP address
		authHeader := c.GetHeader("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userId := validateToken(tokenString, jwtSecret)
			if userId != nil {
				clientIdentifier = *userId
				superbaseClient.AddUsageEvent(superbase.UsageEvent{
					UserId:    *userId,
					EventType: c.Request.URL.Path,
				})
			} else {
				// Invalid token, return unauthorized error
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
		}

		currentTime := time.Now().Unix()

		rateLimitStore.Lock()
		defer rateLimitStore.Unlock()

		// Initialize the request times slice if not present
		if _, exists := rateLimitStore.clients[clientIdentifier]; !exists {
			rateLimitStore.clients[clientIdentifier] = []int64{}
		}

		// Append the current request time
		rateLimitStore.clients[clientIdentifier] = append(rateLimitStore.clients[clientIdentifier], currentTime)

		// Remove timestamps older than the time window
		validTime := currentTime - ratelimitConfig.WindowInSec
		validRequests := []int64{}
		for _, t := range rateLimitStore.clients[clientIdentifier] {
			if t > validTime {
				validRequests = append(validRequests, t)
			}
		}
		rateLimitStore.clients[clientIdentifier] = validRequests

		// Check if the request count exceeds the limit
		if len(validRequests) > ratelimitConfig.RateLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
