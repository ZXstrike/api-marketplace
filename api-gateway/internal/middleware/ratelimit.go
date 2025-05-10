package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const (
	rateLimitPrefix        = "api_rl_config:"
	rateLimitCacheDuration = 5 * time.Minute
	notFoundCacheDuration  = 5 * time.Minute
	notFoundPlaceholder    = "NOT_FOUND"
	userStatusCachePrefix  = "user_status:"
	statusSuspended        = "SUSPENDED"
)

func RateLimitMiddleware(logger *log.Logger, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// apiKeyVal, keyExist := c.Get("apiKey")

		// if !keyExist {
		// 	logger.Println("API Key not found in context")

		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "API Key not found"})
		// 	return
		// }

		c.Next()
	}
}

func getAPIRateLimit(redisClient *redis.Client, apiKey string) (int, error) {
	// Fetch rate limit from Redis
	rateLimit, err := redisClient.Get(rateLimitPrefix + apiKey).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	if rateLimit == "" {
		return 0, nil
	}

	return int(1), nil
}
