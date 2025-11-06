package router

import (
	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/http/handler"
	"github.com/grachmannico95/mileapp-test-be/internal/http/middleware"
	"github.com/grachmannico95/mileapp-test-be/internal/http/router/routes"
)

func NewRouter(cfg *config.Config, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(middleware.SecurityHeadersMiddleware())

	routes.RegisterHealthRoutes(router)

	v1 := router.Group("/api/v1")
	{
		routes.RegisterAuthRoutes(v1, cfg, authHandler)
		routes.RegisterTaskRoutes(v1, cfg, authHandler, taskHandler)
	}

	return router
}
