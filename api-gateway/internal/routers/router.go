package routers

import (
	"net/http"

	"github.com/ZXstrike/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {

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

	apiGroup := router.Group("/api", // Add any middleware here if needed
		gin.Recovery(),
		gin.Logger(),
	)

	apiGroup.Any("/*proxyPath", proxy.ProxyHandler(db, redisClient))
}
