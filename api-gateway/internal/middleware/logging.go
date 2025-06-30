package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoggingMiddleware(logger *log.Logger, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logEntry := fmt.Sprintf("| %3d | %13v | %15s | %-7s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		if errorMessage != "" {
			logEntry += "\n Error: " + errorMessage
		}

		logger.Println(logEntry)

		if statusCode >= 200 && statusCode < 300 {
			apiKey, exists := c.Get("api_key")
			if !exists {
				logger.Println("No API key found in context")
				return
			}

			logEntry := models.UsageLog{
				SubscriptionID:   apiKey.(*models.APIKey).Subscription.ID,
				APIKeyID:         apiKey.(*models.APIKey).ID,
				RequestTimestamp: start,
			}

			if err := db.Create(&logEntry).Error; err != nil {
				logger.Printf("Error logging usage: %v", err)
			}

		} else {
		}
	}
}
