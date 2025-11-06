package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" && isOriginAllowed(origin, cfg.CORS.AllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowedMethods, ", "))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowedHeaders, ", "))

			if len(cfg.CORS.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(cfg.CORS.ExposeHeaders, ", "))
			}

			if cfg.CORS.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			if cfg.CORS.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.CORS.MaxAge))
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}

		if allowed == origin {
			return true
		}

		if strings.HasPrefix(allowed, "*.") {
			domain := allowed[1:] // Remove the *
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}

	return false
}
