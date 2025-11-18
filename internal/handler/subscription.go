package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"subscription-service-go/internal/model"
	"subscription-service-go/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.SubscriptionCreate true "Subscription data"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var createReq model.SubscriptionCreate
	if err := c.ShouldBindJSON(&createReq); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	subscription, err := h.service.CreateSubscription(&createReq)
	if err != nil {
		log.Printf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid subscription ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := h.service.GetSubscription(id)
	if err != nil {
		log.Printf("Failed to get subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
		return
	}

	if subscription == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetSubscriptions godoc
// @Summary Get subscriptions
// @Description Get all subscriptions or filter by user ID
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")

	var subscriptions []model.Subscription
	var err error

	if userID != "" {
		subscriptions, err = h.service.GetUserSubscriptions(userID)
	} else {
		subscriptions, err = h.service.GetAllSubscriptions()
	}

	if err != nil {
		log.Printf("Failed to get subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update an existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body model.SubscriptionUpdate true "Subscription update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid subscription ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var updateReq model.SubscriptionUpdate
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateSubscription(id, &updateReq); err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		log.Printf("Failed to update subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid subscription ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.service.DeleteSubscription(id); err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		log.Printf("Failed to delete subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CalculateCost godoc
// @Summary Calculate subscription cost
// @Description Calculate total cost of subscriptions for a period with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body model.CostRequest true "Cost calculation request"
// @Success 200 {object} model.CostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/cost [post]
func (h *SubscriptionHandler) CalculateCost(c *gin.Context) {
	var costReq model.CostRequest
	if err := c.ShouldBindJSON(&costReq); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	totalCost, err := h.service.CalculateCost(&costReq)
	if err != nil {
		log.Printf("Failed to calculate cost: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate cost"})
		return
	}

	response := model.CostResponse{
		TotalCost:   totalCost,
		Period:      costReq.StartPeriod + " - " + costReq.EndPeriod,
		UserID:      costReq.UserID,
		ServiceName: costReq.ServiceName,
	}

	c.JSON(http.StatusOK, response)
}
