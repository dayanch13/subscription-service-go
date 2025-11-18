package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"subscription-service-go/internal/config"
	"subscription-service-go/internal/handler"
	"subscription-service-go/internal/repository/postgres"
	"subscription-service-go/internal/service"
	"subscription-service-go/pkg/logger"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления онлайн-подписками пользователей
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbConfig := postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := postgres.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository, service, and handler
	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// Initialize Gin router
	router := gin.Default()
	handler.SetupRoutes(router, subscriptionService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
