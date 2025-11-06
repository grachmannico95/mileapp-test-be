package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/dto"
	"github.com/grachmannico95/mileapp-test-be/internal/util"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("authentication required"))
			c.Abort()
			return
		}

		_, err = util.ValidateJWT(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("invalid or expired token"))
			c.Abort()
			return
		}

		c.Next()
	}
}
