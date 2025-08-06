package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		subscriptionID := info.SubscriptionID // Assuming this is available in CachedKeyInfo

		// Phase 1: Reserve funds from the consumer's wallet BEFORE the handler runs.
		err := db.Transaction(func(tx *gorm.DB) error {
			debitTx := tx.Model(&models.Wallet{}).
				Where("user_id = ? AND wallet_type = ? AND balance >= ?", consumerID, "consumer", price).
				Update("balance", gorm.Expr("balance - ?", price))

			if debitTx.Error != nil {
				return debitTx.Error
			}
			if debitTx.RowsAffected == 0 {
				return errors.New("insufficient account balance")
			}

			// Log the debit as a transaction
			debitLog := models.PaymentTransaction{
				UserID:         consumerID,
				SubscriptionID: &subscriptionID,
				Amount:         -price, // Negative amount for debit
				Description:    fmt.Sprintf("Hold for API call on subscription %s", subscriptionID),
			}
			return tx.Create(&debitLog).Error
		})

		if err != nil {
			if err.Error() == "insufficient account balance" {
				c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient account balance"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error during funds reservation"})
			}
			return
		}

		// Defer Phase 2 to run AFTER the handler.
		defer func() {
			if r := recover(); r != nil {
				// Launch refund in a goroutine, but log the panic immediately.
				go logAndRefundConsumer(db, consumerID, subscriptionID, price, "Refund due to server panic")
				panic(r) // Re-panic to allow Gin's recovery middleware to handle it.
			}

			// The logging and crediting logic will now run in the background.
			if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
				// On success, calculate fees and credit provider and platform.
				platformFeePercentage, err := getPlatformFeeSettings(db)
				if err != nil {
					log.Printf("CRITICAL: Could not retrieve platform fee settings: %v. Refunding consumer.", err)
					// Launch refund in a goroutine.
					go logAndRefundConsumer(db, consumerID, subscriptionID, price, "Refund due to internal configuration error")
					return
				}

				platformCut := price * platformFeePercentage
				providerCut := price - platformCut

				// Launch these database operations asynchronously.
				go logAndCreditProvider(db, providerID, subscriptionID, providerCut)
				go logAndCreditPlatform(db, subscriptionID, platformCut)
			} else {
				// On failure, refund the consumer asynchronously.
				description := fmt.Sprintf("Refund for failed API call (HTTP Status %d)", c.Writer.Status())
				go logAndRefundConsumer(db, consumerID, subscriptionID, price, description)
			}
		}()

		// Proceed to the actual request handler without holding a transaction.
		c.Next()
	}
}

func getPlatformFeeSettings(db *gorm.DB) (float64, error) {
	var setting models.SystemSetting
	// Fetch the specific setting for platform fee percentage.
	if err := db.Where("key = ?", "platform_fee_percentage").First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("system setting 'platform_fee_percentage' not found")
		}
		return 0, err
	}

	// Parse the string value from the database into a float.
	fee, err := strconv.ParseFloat(setting.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse platform fee percentage value '%s': %w", setting.Value, err)
	}

	// The fee should be a decimal (e.g., 0.05 for 5%).
	return fee, nil
}

// logAndCreditProvider creates a transaction log and credits the provider's wallet.
func logAndCreditProvider(db *gorm.DB, providerID, subscriptionID string, amount float64) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Create PaymentTransaction record for the credit.
		creditLog := models.PaymentTransaction{
			UserID:         providerID,
			SubscriptionID: &subscriptionID,
			Amount:         amount,
			Description:    fmt.Sprintf("Revenue from API call on subscription %s", subscriptionID),
		}
		if err := tx.Create(&creditLog).Error; err != nil {
			return err
		}

		// 2. Update the provider's wallet.
		return tx.Model(&models.Wallet{}).Where("user_id = ? AND wallet_type = ?", providerID, "provider").Update("balance", gorm.Expr("balance + ?", amount)).Error
	})

	if err != nil {
		log.Printf("CRITICAL: FAILED TO PROCESS PROVIDER CREDIT for %s. Error: %v", providerID, err)
	}
}

// logAndCreditPlatform creates a PlatformRevenue record and credits the platform's wallet.
func logAndCreditPlatform(db *gorm.DB, subscriptionID string, amount float64) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Create a PlatformRevenue record for the fee.
		revenueLog := models.PlatformRevenue{
			// SourcePayoutID is nil because this is a direct per-call fee, not from a payout.
			Amount:      amount,
			Description: fmt.Sprintf("Platform fee from API call on subscription %s", subscriptionID),
		}
		if err := tx.Create(&revenueLog).Error; err != nil {
			return err
		}

		// 2. Update the platform's main wallet.
		return tx.Model(&models.PlatformWallet{}).Where("name = ?", "main").Update("balance", gorm.Expr("balance + ?", amount)).Error
	})

	if err != nil {
		log.Printf("CRITICAL: FAILED TO PROCESS PLATFORM FEE. Error: %v", err)
	}
}

// logAndRefundConsumer creates a transaction log and refunds the consumer's wallet.
func logAndRefundConsumer(db *gorm.DB, consumerID, subscriptionID string, amount float64, reason string) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Create PaymentTransaction record for the refund.
		refundLog := models.PaymentTransaction{
			UserID:         consumerID,
			SubscriptionID: &subscriptionID,
			Amount:         amount, // Positive amount for credit/refund
			Description:    reason,
		}
		if err := tx.Create(&refundLog).Error; err != nil {
			return err
		}

		// 2. Update the consumer's wallet.
		return tx.Model(&models.Wallet{}).Where("user_id = ? AND wallet_type = ?", consumerID, "consumer").Update("balance", gorm.Expr("balance + ?", amount)).Error
	})

	if err != nil {
		log.Printf("CRITICAL: FAILED TO PROCESS CONSUMER REFUND for %s. Error: %v", consumerID, err)
	}
}
