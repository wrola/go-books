package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// APIKeyHeader is the header name for the API key
	APIKeyHeader = "X-API-Key"
)

// AuthMiddleware returns a Gin middleware for API key authentication
func AuthMiddleware() gin.HandlerFunc {
	// Get API key from environment
	apiKey := os.Getenv("API_KEY")

	// Check if auth is disabled (for development)
	authDisabled := os.Getenv("AUTH_DISABLED") == "true"

	return func(c *gin.Context) {
		// Skip auth if disabled
		if authDisabled {
			c.Next()
			return
		}

		// Skip auth if no API key is configured (allow open access)
		if apiKey == "" {
			c.Next()
			return
		}

		// Get API key from header
		providedKey := c.GetHeader(APIKeyHeader)

		// Check if key is provided
		if providedKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
			})
			return
		}

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(providedKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware returns middleware that extracts API key but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey != "" {
			c.Set("api_key", apiKey)
		}
		c.Next()
	}
}

// SkipAuthPaths returns middleware that skips authentication for specific paths
func SkipAuthPaths(paths ...string) gin.HandlerFunc {
	pathSet := make(map[string]bool)
	for _, p := range paths {
		pathSet[p] = true
	}

	return func(c *gin.Context) {
		// Check if current path should skip auth
		path := c.Request.URL.Path
		for p := range pathSet {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		// Apply auth middleware
		AuthMiddleware()(c)
	}
}
