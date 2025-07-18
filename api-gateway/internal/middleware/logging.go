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

		// Only log usage for successful requests.
		if statusCode >= 200 && statusCode < 300 {
			// Use the lightweight CachedKeyInfo object from the context.
			infoVal, exists := c.Get("api_key")
			if !exists {
				logger.Println("No api_key found in context for usage logging")
				return
			}

			info := infoVal.(*CachedKeyInfo)

			// Perform the database write in a separate goroutine.
			go func() {
				usageLog := models.UsageLog{
					// The required IDs are now directly on the cached info struct.
					SubscriptionID:   info.APIID, // Assuming APIID in CachedKeyInfo corresponds to SubscriptionID's purpose
					APIKeyID:         info.APIKeyID,
					RequestTimestamp: start,
				}

				if err := db.Create(&usageLog).Error; err != nil {
					// Use the same logger to report the background task failure.
					logger.Printf("Error logging usage to DB: %v", err)
				}
			}()
		}
	}
}
