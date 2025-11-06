package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/dto"
	"github.com/grachmannico95/mileapp-test-be/internal/util"
)

const (
	CSRFTokenHeader = "X-CSRF-Token"
)

func CSRFMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		headerToken := c.GetHeader(CSRFTokenHeader)
		if headerToken == "" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse("csrf token missing in header"))
			c.Abort()
			return
		}

		cookieToken, err := c.Cookie("csrf_token")
		if err != nil {
			c.JSON(http.StatusForbidden, dto.ErrorResponse("csrf token missing in cookie"))
			c.Abort()
			return
		}

		if headerToken != cookieToken {
			c.JSON(http.StatusForbidden, dto.ErrorResponse("csrf token mismatch"))
			c.Abort()
			return
		}

		if !util.ValidateCSRFToken(headerToken, cfg.CSRF.Secret) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse("invalid csrf token"))
			c.Abort()
			return
		}

		c.Next()
	}
}
