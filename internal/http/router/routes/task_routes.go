package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/http/handler"
	"github.com/grachmannico95/mileapp-test-be/internal/http/middleware"
)

func RegisterTaskRoutes(v1 *gin.RouterGroup, cfg *config.Config, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler) {
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	protected.Use(middleware.CSRFMiddleware(cfg))
	{
		protected.GET("/tasks", taskHandler.List)
		protected.GET("/tasks/:id", taskHandler.GetByID)
		protected.POST("/tasks", taskHandler.Create)
		protected.PUT("/tasks/:id", taskHandler.Update)
		protected.DELETE("/tasks/:id", taskHandler.Delete)
	}
}
