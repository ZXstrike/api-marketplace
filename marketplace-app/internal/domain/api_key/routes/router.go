package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/api_key/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/api_key/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/api_key/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	service := service.New(repo, privateKey, publicKey)
	h := handler.New(service)

	apiKey := r.Group("/api-keys")
	{
		apiKey.POST("/create", middleware.AuthMiddleware(publicKey), h.CreateAPIKey)
		apiKey.DELETE("/delete", middleware.AuthMiddleware(publicKey), h.DeleteAPIKey)
	}

}
