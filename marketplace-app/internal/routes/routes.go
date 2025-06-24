package routes

import (
	"crypto/ecdsa"

	apiRoutes "github.com/ZXstrike/marketplace-app/internal/domain/api/routes"
	apiKeyRoutes "github.com/ZXstrike/marketplace-app/internal/domain/api_key/routes"
	authRoutes "github.com/ZXstrike/marketplace-app/internal/domain/auth/routes"
	storeRoutes "github.com/ZXstrike/marketplace-app/internal/domain/store/routes"
	subscriptionRoutes "github.com/ZXstrike/marketplace-app/internal/domain/subscription/routes"
	userRoutes "github.com/ZXstrike/marketplace-app/internal/domain/user/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {

	router.Static("/files", "./files")

	// Auth routes
	authRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// User routes
	userRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// Store routes
	storeRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// API routes
	apiRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// Subscription routes
	subscriptionRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)

	// Api Key routes
	apiKeyRoutes.RegisterRoutes(&router.RouterGroup, db, privateKey, publicKey)
}
