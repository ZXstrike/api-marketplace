package handler

import (
	"net/http"

	"github.com/ZXstrike/marketplace-app/internal/domain/billing/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetBillingInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletType := c.Query("wallet_type")
	if walletType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_type query parameter is required"})
		return
	}

	billingInfo, err := h.service.GetBillingInfo(userID.(string), walletType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve billing info"})
		return
	}

	c.JSON(http.StatusOK, billingInfo)
}

func (h *Handler) TopUp(c *gin.Context) {
	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Top-up is only allowed for the "consumer" wallet.
	if err := h.service.TopUp(userID.(string), "consumer", request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance topped up successfully"})
}

func (h *Handler) ProcessPayment(c *gin.Context) {
	var request struct {
		Amount     float64 `json:"amount" binding:"required,gt=0"`
		WalletType string  `json:"wallet_type" binding:"required"`
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := h.service.ProcessPayment(userID.(string), request.WalletType, request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}

func (h *Handler) GetPaymentHistory(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	history, err := h.service.GetPaymentHistory(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_history": history,
	})
}
