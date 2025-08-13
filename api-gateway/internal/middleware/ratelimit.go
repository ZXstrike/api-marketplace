package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9" // Correct import path for v9
	"gorm.io/gorm"
)

const (
	userStatusCachePrefix = "user_status:"
	statusSuspended       = "SUSPENDED"
	rateLimitPrefix       = "rate_limit:"
	defaultRateLimit      = 60 // Default requests per minute if DB lookup fails
	rateLimitWindow       = 1 * time.Minute
	rateLimitSettingKey   = "rate_limit_per_minute"
	settingCacheTTL       = 5 * time.Minute
)

// Global cache for the rate limit setting to minimize DB queries.
var (
	cachedRateLimit int
	lastFetch       time.Time
	mu              sync.RWMutex
)

// RateLimitMiddleware enforces a request rate limit based on a setting from the database.
// It uses Redis for counting requests and has an in-memory cache for the rate limit setting.
func RateLimitMiddleware(logger *log.Logger, redisClient *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context() // Use context for Redis operations

		infoVal, keyExist := c.Get("api_key")
		if !keyExist {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "API Key not found"})
			return
		}
		info := infoVal.(*CachedKeyInfo)

		apiKey := info.APIKeyID
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "API Key is required"})
			return
		}

		// Check if the user is suspended
		userStatus, err := redisClient.Get(ctx, userStatusCachePrefix+info.ConsumerID).Result()
		if err != nil && err != redis.Nil {
			logger.Printf("Error checking user status: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if userStatus == statusSuspended {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is suspended"})
			return
		}

		// Get the rate limit, fetching from DB via local cache if necessary
		limit := getRateLimit(ctx, db, logger)
		key := rateLimitPrefix + apiKey

		// Use a Redis pipeline to make INCR and EXPIRE atomic
		pipe := redisClient.Pipeline()
		incrCmd := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, rateLimitWindow)
		_, err = pipe.Exec(ctx) // Pass context to Exec
		if err != nil {
			logger.Printf("Error executing rate limit pipeline: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not process request"})
			return
		}

		current := incrCmd.Val()

		// Set informative headers
		remaining := limit - int(current)
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(int64(rateLimitWindow.Seconds()), 10))

		// Check if the limit is exceeded
		if current > int64(limit) {
			logger.Printf("Rate limit exceeded for key: %s", apiKey)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}

// getRateLimit fetches the rate limit from a local in-memory cache or the database.
// The cache helps to avoid hitting the database on every single request.
func getRateLimit(ctx context.Context, db *gorm.DB, logger *log.Logger) int {
	mu.RLock()
	// Check if the cached value is still fresh
	if time.Since(lastFetch) < settingCacheTTL {
		limit := cachedRateLimit
		mu.RUnlock()
		return limit
	}
	mu.RUnlock()

	// If cache is stale, acquire a write lock to update it
	mu.Lock()
	defer mu.Unlock()

	// Double-check if another request already updated the cache while we were waiting for the lock
	if time.Since(lastFetch) < settingCacheTTL {
		return cachedRateLimit
	}

	// Fetch the setting from the database
	var setting models.SystemSetting
	if err := db.WithContext(ctx).Where("key = ?", rateLimitSettingKey).First(&setting).Error; err != nil {
		logger.Printf("Could not fetch rate limit setting, using default (%d): %v", defaultRateLimit, err)
		cachedRateLimit = defaultRateLimit
	} else {
		limit, err := strconv.Atoi(setting.Value)
		if err != nil {
			logger.Printf("Invalid rate limit value in database ('%s'), using default (%d): %v", setting.Value, defaultRateLimit, err)
			cachedRateLimit = defaultRateLimit
		} else {
			cachedRateLimit = limit
		}
	}

	lastFetch = time.Now()
	return cachedRateLimit
}
