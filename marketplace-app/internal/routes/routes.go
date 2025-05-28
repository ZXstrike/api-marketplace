package routes

import (
	"crypto/ecdsa"

	authRoutes "github.com/ZXstrike/marketplace-app/internal/domain/auth/routes"
	storeRoutes "github.com/ZXstrike/marketplace-app/internal/domain/store/routes"
	userRoutes "github.com/ZXstrike/marketplace-app/internal/domain/user/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {

	// Auth routes
	authRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// User routes
	userRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// Store routes
	storeRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)
}
