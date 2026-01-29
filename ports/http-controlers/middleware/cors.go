package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a Gin middleware for CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	// Get allowed origins from environment, default to none (strict)
	allowedOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	// Get allowed methods from environment, default to common methods
	allowedMethodsStr := os.Getenv("CORS_ALLOWED_METHODS")
	allowedMethods := "GET, POST, PUT, DELETE, OPTIONS"
	if allowedMethodsStr != "" {
		allowedMethods = allowedMethodsStr
	}

	// Get allowed headers from environment
	allowedHeadersStr := os.Getenv("CORS_ALLOWED_HEADERS")
	allowedHeaders := "Content-Type, Authorization, X-API-Key, X-Request-ID"
	if allowedHeadersStr != "" {
		allowedHeaders = allowedHeadersStr
	}

	// Get exposed headers from environment
	exposedHeadersStr := os.Getenv("CORS_EXPOSED_HEADERS")
	exposedHeaders := "X-Request-ID"
	if exposedHeadersStr != "" {
		exposedHeaders = exposedHeadersStr
	}

	// Check if credentials are allowed
	allowCredentials := os.Getenv("CORS_ALLOW_CREDENTIALS") == "true"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		originAllowed := false
		for _, allowed := range allowedOrigins {
			if allowed == "*" || allowed == origin {
				originAllowed = true
				break
			}
		}

		if originAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", allowedMethods)
			c.Header("Access-Control-Allow-Headers", allowedHeaders)
			c.Header("Access-Control-Expose-Headers", exposedHeaders)
			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			c.Header("Access-Control-Max-Age", "86400") // 24 hours
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
