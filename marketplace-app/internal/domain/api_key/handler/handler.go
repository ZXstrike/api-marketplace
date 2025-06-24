package handler

import (
	"github.com/ZXstrike/marketplace-app/internal/domain/api_key/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateAPIKey(c *gin.Context) {
	var req struct {
		SubscriptionID string `json:"subscription_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	apiKey, err := h.service.CreateAPIKey(req.SubscriptionID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(200, gin.H{"api_key": apiKey})
	return
}

func (h *Handler) DeleteAPIKey(c *gin.Context) {
	apiKeyID := c.Query("api_key_id")

	err := h.service.DeleteAPIKey(apiKeyID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.JSON(200, gin.H{"message": "API key deleted successfully"})
}
