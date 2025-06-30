package middleware

import (
	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BillingMiddleware(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		// This middleware can be used to check billing status, subscription validity, etc.
		// For now, it just passes the request through.

		// You can implement your billing logic here, such as checking if the user's subscription is active,
		// if they have sufficient balance, etc.

		c.Next() // Proceed to the next middleware or handler

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			apiKey, exists := c.Get("api_key")
			if !exists {
				c.JSON(500, gin.H{"error": "API key not found in context"})
				return
			}

			user := apiKey.(*models.APIKey).Subscription.Consumer
			if user.AccountBalance < apiKey.(*models.APIKey).Subscription.APIVersion.PricePerCall {
				c.JSON(402, gin.H{"error": "Insufficient account balance"})
				return
			}

			user.AccountBalance -= apiKey.(*models.APIKey).Subscription.APIVersion.PricePerCall
			if err := db.Save(user).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to update account balance"})
				return
			}

		}
	}

}
