package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/subscription/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/subscription/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/subscription/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	service := service.New(repo, privateKey, publicKey)
	h := handler.New(service)

	substription := r.Group("/subscriptions")
	{
		substription.POST("/subscribe", middleware.AuthMiddleware(publicKey), h.SubscribeToAPI)
		substription.POST("/unsubscribe", middleware.AuthMiddleware(publicKey), h.UnsubscribeFromAPI)
		substription.GET("/get", middleware.AuthMiddleware(publicKey), h.GetSubscription)
		substription.GET("/get-by-user", middleware.AuthMiddleware(publicKey), h.GetSubscriptions)
	}
}
