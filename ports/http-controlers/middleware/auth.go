package middleware

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	APIKeyHeader = "X-API-Key"
)

func AuthMiddleware() gin.HandlerFunc {
	apiKey := os.Getenv("API_KEY")

	// Check if auth is explicitly disabled (for development only)
	// Requires AUTH_DISABLED=true - empty API_KEY alone does NOT disable auth
	authDisabled := os.Getenv("AUTH_DISABLED") == "true"

	if authDisabled {
		log.Println("WARNING: Authentication is disabled (AUTH_DISABLED=true). Do not use in production.")
	}

	return func(c *gin.Context) {
		if authDisabled {
			c.Next()
			return
		}

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "server configuration error",
			})
			log.Println("ERROR: API_KEY not configured. Set API_KEY or AUTH_DISABLED=true")
			return
		}

		providedKey := c.GetHeader(APIKeyHeader)

		if providedKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
			})
			return
		}

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(providedKey)) != 1 {
			log.Printf("AUTH_FAILURE: Invalid API key from IP %s", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			return
		}

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey != "" {
			c.Set("api_key", apiKey)
		}
		c.Next()
	}
}

func SkipAuthPaths(paths ...string) gin.HandlerFunc {
	pathSet := make(map[string]bool)
	for _, p := range paths {
		pathSet[p] = true
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for p := range pathSet {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		AuthMiddleware()(c)
	}
}
