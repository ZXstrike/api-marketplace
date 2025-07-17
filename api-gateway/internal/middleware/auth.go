package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CachedKeyInfo is a lightweight struct for caching essential API key data.
type CachedKeyInfo struct {
	APIKeyID     string  `json:"api_key_id"`
	ConsumerID   string  `json:"consumer_id"`
	ProviderID   string  `json:"provider_id"`
	APIID        string  `json:"api_id"`
	PricePerCall float64 `json:"price_per_call"`
	Balance      float64 `json:"balance"`
}

const (
	apiKeyCacheDuration    = 5 * time.Minute
	apiKeyNotFoundDuration = 1 * time.Minute
	apiKeyCachePrefix      = "apikey:"
	apiKeyNotFoundValue    = "not_found"
)

func AuthMiddleware(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyString := c.GetHeader("Api-Key")
		if apiKeyString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		// 1. Check cache first
		cachedInfo, err := getAPIKeyFromCache(c.Request.Context(), redisClient, apiKeyString)
		if err == nil && cachedInfo != nil {
			// Cache hit
		} else if err != nil && err.Error() == apiKeyNotFoundValue {
			// "Not found" cached, treat as invalid
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		} else {
			// Cache miss or other Redis error, fallback to DB
			cachedInfo, err = ValidateAPIKey(db, apiKeyString)
			if err != nil {
				// Cache the "not found" result to prevent DB hammering
				if err.Error() == "invalid API key" {
					cacheAPIKey(c.Request.Context(), redisClient, apiKeyString, nil)
				}
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
				c.Abort()
				return
			}
			// DB lookup successful, cache the result
			cacheAPIKey(c.Request.Context(), redisClient, apiKeyString, cachedInfo)
		}

		// Use the lightweight cached object instead of the heavy GORM model
		c.Set("api_key_info", cachedInfo) // Pass the whole info struct
		c.Set("user_id", cachedInfo.ConsumerID)
		c.Set("api_id", cachedInfo.APIID)

		if cachedInfo.Balance < cachedInfo.PricePerCall {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient account balance"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getAPIKeyFromCache retrieves lightweight key info from Redis cache.
func getAPIKeyFromCache(ctx context.Context, redisClient *redis.Client, keyString string) (*CachedKeyInfo, error) {
	cacheKey := apiKeyCachePrefix + keyString
	cachedData, err := redisClient.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err // Redis error
	}

	if cachedData == apiKeyNotFoundValue {
		return nil, errors.New(apiKeyNotFoundValue) // "Not found" placeholder hit
	}

	var info CachedKeyInfo
	if err := json.Unmarshal([]byte(cachedData), &info); err != nil {
		return nil, err // Deserialization error
	}

	return &info, nil
}

// cacheAPIKey stores lightweight key info in Redis cache.
func cacheAPIKey(ctx context.Context, redisClient *redis.Client, keyString string, info *CachedKeyInfo) {
	cacheKey := apiKeyCachePrefix + keyString

	if info == nil {
		// Cache the "not found" status
		redisClient.Set(ctx, cacheKey, apiKeyNotFoundValue, apiKeyNotFoundDuration)
		return
	}

	// Serialize and cache the found key info
	serializedData, err := json.Marshal(info)
	if err != nil {
		// Log marshalling error, but don't block the request
		return
	}
	redisClient.Set(ctx, cacheKey, serializedData, apiKeyCacheDuration)
}

// ValidateAPIKey finds a key and its related data in a single query.
// It returns a lightweight CachedKeyInfo object if found, otherwise an error.
func ValidateAPIKey(db *gorm.DB, keyString string) (*CachedKeyInfo, error) {
	if keyString == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	incomingKeyHash := sha256.Sum256([]byte(keyString))

	var result CachedKeyInfo
	err := db.Model(&models.APIKey{}).
		Select(
			"api_keys.id as api_key_id",
			"consumers.id as consumer_id",
			"providers.id as provider_id",
			"apis.id as api_id",
			"api_versions.price_per_call",
			"consumers.account_balance as balance",
		).
		Joins("JOIN subscriptions ON subscriptions.id = api_keys.subscription_id").
		Joins("JOIN users as consumers ON consumers.id = subscriptions.consumer_id").
		Joins("JOIN api_versions ON api_versions.id = subscriptions.api_version_id").
		Joins("JOIN apis ON apis.id = api_versions.api_id").
		Joins("JOIN users as providers ON providers.id = apis.provider_id").
		Where("api_keys.key_value_hash = ?", incomingKeyHash[:]).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid API key")
		}
		return nil, fmt.Errorf("database error validating key: %w", err)
	}

	return &result, nil
}
