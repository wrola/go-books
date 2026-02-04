package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	allowedOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	allowedMethodsStr := os.Getenv("CORS_ALLOWED_METHODS")
	allowedMethods := "GET, POST, PUT, DELETE, OPTIONS"
	if allowedMethodsStr != "" {
		allowedMethods = allowedMethodsStr
	}

	allowedHeadersStr := os.Getenv("CORS_ALLOWED_HEADERS")
	allowedHeaders := "Content-Type, Authorization, X-API-Key, X-Request-ID"
	if allowedHeadersStr != "" {
		allowedHeaders = allowedHeadersStr
	}

	exposedHeadersStr := os.Getenv("CORS_EXPOSED_HEADERS")
	exposedHeaders := "X-Request-ID"
	if exposedHeadersStr != "" {
		exposedHeaders = exposedHeadersStr
	}

	allowCredentials := os.Getenv("CORS_ALLOW_CREDENTIALS") == "true"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		originAllowed := false
		hasWildcard := false
		for _, allowed := range allowedOrigins {
			if allowed == "*" {
				hasWildcard = true
			}
			if allowed == "*" || allowed == origin {
				originAllowed = true
				break
			}
		}

		if originAllowed {
			if hasWildcard && allowCredentials {
				if origin != "" {
					c.Header("Access-Control-Allow-Origin", origin)
				}
			} else if hasWildcard {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
			}
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
