package handler

import (
	"github.com/ZXstrike/marketplace-app/internal/domain/subscription/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) SubscribeToAPI(c *gin.Context) {
	var req struct {
		APIVersionID string `json:"api_version_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := h.service.SubscribeToAPI(userID.(string), req.APIVersionID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to subscribe to API"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully subscribed to API"})
}

func (h *Handler) UnsubscribeFromAPI(c *gin.Context) {
	var req struct {
		SubscriptionID string `json:"subscription_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := h.service.UnsubscribeFromAPI(userID.(string), req.SubscriptionID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to unsubscribe from API"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully unsubscribed from API"})
}

func (h *Handler) GetSubscription(c *gin.Context) {
	subscriptionID := c.Query("subscriptionID")

	if subscriptionID == "" {
		c.JSON(400, gin.H{"error": "API Version ID is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	subscription, err := h.service.GetSubscription(userID.(string), subscriptionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve subscription"})
		return
	}

	c.JSON(200, subscription)
}

func (h *Handler) GetSubscriptions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	subscriptions, err := h.service.GetSubscriptionsByUserID(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	c.JSON(200, subscriptions)
}
