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

	billingInfo, err := h.service.GetBillingInfo(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve billing info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"billing_info": billingInfo,
	})

}

func (h *Handler) UpdateBalance(c *gin.Context) {
	// This method will handle the request to update the balance of a user.
	// Implementation will depend on the specific requirements and data structure.
	var request struct {
		Amount float64 `json:"amount"`
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.service.UpdateBalanceByUserID(userID.(string), request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance updated successfully"})
}

func (h *Handler) ProcessPayment(c *gin.Context) {
	// This method will handle the request to process a payment.
	// Implementation will depend on the specific requirements and payment gateway used.
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
