package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/auth/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/auth/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/auth/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	service := service.New(repo, privateKey, publicKey)
	h := handler.New(service)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.RegisterHandler)
		auth.POST("/login", h.LoginHandler)
		auth.POST("/refresh", middleware.AuthMiddleware(publicKey), h.RefreshHandler)
	}
}
