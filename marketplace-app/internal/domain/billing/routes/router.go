package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/billing/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/billing/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/billing/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	service := service.New(repo, privateKey, publicKey)
	h := handler.New(service)

	billing := r.Group("/billing")
	{
		billing.GET("/info", middleware.AuthMiddleware(publicKey), h.GetBillingInfo)
		billing.PUT("/update-balance", middleware.AuthMiddleware(publicKey), h.TopUp)
		billing.POST("/payment", middleware.AuthMiddleware(publicKey), h.ProcessPayment)
		billing.GET("/history", middleware.AuthMiddleware(publicKey), h.GetPaymentHistory)
		billing.POST("/payout", middleware.AuthMiddleware(publicKey), h.HandleCreatePayout)
	}

}
