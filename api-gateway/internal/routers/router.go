package routers

import (
	"log"
	"net/http"

	"github.com/ZXstrike/api-gateway/internal/middleware"
	"github.com/ZXstrike/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Check database connection
		if err := db.Exec("SELECT 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "Database connection failed",
			})
			return
		}
		// Check Redis connection
		if err := redisClient.Ping(c.Request.Context()).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Use NoRoute to handle all unmatched routes as proxy routes
	// This is the proper way to handle catch-all routing in Gin
	router.NoRoute(
		middleware.AuthMiddleware(db, redisClient),
		middleware.LoggingMiddleware(
			log.New(
				log.Writer(),
				"api-gateway: ",
				log.LstdFlags|log.Lshortfile,
			),
			db,
		),
		middleware.BillingMiddleware(db),
		proxy.ProxyHandler(db, redisClient),
	)
}
