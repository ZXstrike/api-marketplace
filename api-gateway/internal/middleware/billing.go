package middleware

import (
	"log"
	"net/http"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BillingMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Correctly get "api_key_info" from the context set by AuthMiddleware.
		infoVal, exists := c.Get("api_key")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "API key info not found in context for billing"})
			return
		}

		// Cast to the correct lightweight struct type.
		info := infoVal.(*CachedKeyInfo)
		price := info.PricePerCall
		consumerID := info.ConsumerID
		providerID := info.ProviderID

		// Phase 1: Reserve funds from the consumer BEFORE the handler runs.
		// This is a very short, single-statement operation.
		debitTx := db.Model(&models.User{}).
			Where("id = ? AND account_balance >= ?", consumerID, price).
			Update("account_balance", gorm.Expr("account_balance - ?", price))

		if debitTx.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error during funds reservation"})
			return
		}

		if debitTx.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient account balance"})
			return
		}

		// Defer Phase 2 to run AFTER the handler.
		defer func() {
			if r := recover(); r != nil {
				refundConsumer(db, consumerID, price)
				panic(r)
			}

			if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
				creditProvider(db, providerID, price)
			} else {
				refundConsumer(db, consumerID, price)
			}
		}()

		// Proceed to the actual request handler without holding a transaction.
		c.Next()
	}
}

// creditProvider credits the provider's account.
func creditProvider(db *gorm.DB, providerID string, price float64) {
	if err := db.Model(&models.User{}).Where("id = ?", providerID).Update("account_balance", gorm.Expr("account_balance + ?", price)).Error; err != nil {
		log.Printf("CRITICAL: FAILED TO CREDIT PROVIDER %s. Error: %v", providerID, err)
	}
}

// refundConsumer refunds the consumer if the request fails after debit.
func refundConsumer(db *gorm.DB, consumerID string, price float64) {
	if err := db.Model(&models.User{}).Where("id = ?", consumerID).Update("account_balance", gorm.Expr("account_balance + ?", price)).Error; err != nil {
		log.Printf("CRITICAL: FAILED TO REFUND CONSUMER %s. Error: %v", consumerID, err)
	}
}
