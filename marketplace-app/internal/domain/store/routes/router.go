package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/store/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/store/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/store/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	service := service.New(repo, privateKey, publicKey)
	h := handler.New(service)

	store := r.Group("/store")
	{
		store.GET("/user/:userID", h.GetStoreByUserIDHandler)
		store.GET("/username/:username", h.GetStoreByUsernameHandler)
		store.GET("/all", h.GetAllStoresHandler)
		store.POST("/create", middleware.AuthMiddleware(publicKey), h.CreateStoreHandler)
		store.PUT("/update", middleware.AuthMiddleware(publicKey), h.UpdateStoreHandler)
		store.GET("/apis", middleware.AuthMiddleware(publicKey), h.GetStoreApisHandler)
	}
}
