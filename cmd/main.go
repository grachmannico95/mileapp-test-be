package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/http/handler"
	"github.com/grachmannico95/mileapp-test-be/internal/http/router"
	"github.com/grachmannico95/mileapp-test-be/internal/http/server"
	"github.com/grachmannico95/mileapp-test-be/internal/repository"
	"github.com/grachmannico95/mileapp-test-be/internal/service"
	"github.com/grachmannico95/mileapp-test-be/pkg/database"
)

func main() {
	// init config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	gin.SetMode(cfg.Server.GinMode)

	// init database connection
	mongoDB, err := database.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database, cfg.MongoDB.Timeout)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoDB.Close(ctx); err != nil {
			log.Printf("error closing mongodb connection: %v", err)
		}
	}()

	log.Println("successfully connected to mongodb")

	// setup database indexes
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := database.SetupIndexes(ctx, mongoDB.Database); err != nil {
		log.Fatalf("failed to setup database indexes: %v", err)
	}
	log.Println("database indexes created successfully")

	// inject repositories
	userRepo := repository.NewUserRepository(mongoDB.Database)
	taskRepo := repository.NewTaskRepository(mongoDB.Database)

	// inject services
	authService := service.NewAuthService(userRepo, cfg)
	taskService := service.NewTaskService(taskRepo)

	// inject handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	taskHandler := handler.NewTaskHandler(taskService)

	// init router
	r := router.NewRouter(cfg, authHandler, taskHandler)

	// init server
	srv := server.NewServer(cfg.Server.Port, r)
	log.Printf("starting server on :%s", cfg.Server.Port)

	// start server in goroutine
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.GetHTTPServer().Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
