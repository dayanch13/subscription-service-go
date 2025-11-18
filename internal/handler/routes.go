package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"subscription-service-go/internal/service"
)

func SetupRoutes(router *gin.Engine, subscriptionService *service.SubscriptionService) {
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Subscription routes
	api := router.Group("/api/v1")
	{
		api.POST("/subscriptions", subscriptionHandler.CreateSubscription)
		api.GET("/subscriptions/:id", subscriptionHandler.GetSubscription)
		api.GET("/subscriptions", subscriptionHandler.GetSubscriptions)
		api.PUT("/subscriptions/:id", subscriptionHandler.UpdateSubscription)
		api.DELETE("/subscriptions/:id", subscriptionHandler.DeleteSubscription)
		api.POST("/subscriptions/cost", subscriptionHandler.CalculateCost)
	}
}
