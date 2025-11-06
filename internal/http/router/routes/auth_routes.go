package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/http/handler"
	"github.com/grachmannico95/mileapp-test-be/internal/http/middleware"
)

func RegisterAuthRoutes(v1 *gin.RouterGroup, cfg *config.Config, authHandler *handler.AuthHandler) {
	v1.POST("/login", authHandler.Login)
	v1.POST("/logout", authHandler.Logout, middleware.AuthMiddleware(cfg), middleware.CSRFMiddleware(cfg))
}
