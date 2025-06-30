package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Check for a specific header or apiKey

		apiKey := c.GetHeader("Api-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop the request if unauthorized
			return
		}

		apiKeyModel, err := ValidateAPIKey(db, apiKey)
		if err != nil {
			// If the key is invalid, return an error response
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort() // Stop the request if the key is invalid
			return
		}

		// If the key is valid, you can set it in the context for later use
		c.Set("api_key", apiKeyModel)
		// Optionally, you can also set the user ID or other relevant information
		c.Set("user_id", apiKeyModel.Subscription.Consumer.ID)
		// You can also set the API ID if needed
		c.Set("api_id", apiKeyModel.Subscription.APIVersion.API.ID)

		userBalance := apiKeyModel.Subscription.Consumer.AccountBalance
		apiPrice := apiKeyModel.Subscription.APIVersion.PricePerCall

		if userBalance < apiPrice {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient account balance"})
			c.Abort() // Stop the request if the balance is insufficient
			return
		}

		c.Next()
	}
}

// ValidateAPIKey finds a key in the database based on the provided plaintext string.
// It returns the valid APIKey object if found, otherwise an error.
func ValidateAPIKey(db *gorm.DB, keyString string) (*models.APIKey, error) {
	// 1. Ensure the key string is not empty.
	if keyString == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}
	// 2. Split the key to get the prefix. The format is assumed to be "prefix_randompart".
	// Since the prefix itself can contain underscores (e.g., "mk_live_demo"),
	// we split at the last underscore to separate the prefix from the random part.
	lastUnderscoreIndex := strings.LastIndex(keyString, "_")
	if lastUnderscoreIndex == -1 {
		return nil, fmt.Errorf("invalid API key format: missing separator")
	}
	prefix := keyString[:lastUnderscoreIndex]

	_ = prefix // Use the prefix if needed, e.g., for logging or validation
	// 3. Hash the entire incoming key string to match what's in the database.
	incomingKeyHash := sha256.Sum256([]byte(keyString))

	// 4. Query the database for the key.
	// This query is highly efficient as it can use a composite index
	// on (key_prefix, key_value_hash) or individual indexes on both.
	var apiKey models.APIKey
	err := db.Where("key_value_hash = ?", incomingKeyHash[:]).
		Preload("Subscription").
		Preload("Subscription.Consumer").
		Preload("Subscription.APIVersion").
		Preload("Subscription.APIVersion.API").
		First(&apiKey).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// This is the "invalid key" case.
			return nil, fmt.Errorf("invalid API key")
		}
		// This is a different, unexpected database error.
		return nil, fmt.Errorf("database error validating key: %w", err)
	}

	// 5. Key is valid and found. Return the model instance.
	return &apiKey, nil
}
