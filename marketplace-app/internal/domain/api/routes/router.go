package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/api/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/api/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/api/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	service := service.New(repo, privateKey, publicKey)
	h := handler.New(service)

	api := r.Group("/api")
	{
		api.POST("/create", middleware.AuthMiddleware(publicKey), h.CreateNewAPI)
		api.PUT("/update", middleware.AuthMiddleware(publicKey), h.UpdateAPI)
		api.DELETE("/delete/:id", middleware.AuthMiddleware(publicKey), h.DeleteAPI)
		api.GET("/all", h.GetAllAPIs)
		api.GET("/:id", h.GetAPIByID)
		api.POST("/create-endpoint", middleware.AuthMiddleware(publicKey), h.CreateNewAPIEndpoint)
		api.PUT("/update-endpoint", middleware.AuthMiddleware(publicKey), h.UpdateAPIEndpoint)
		api.DELETE("/delete-endpoint/:id", middleware.AuthMiddleware(publicKey), h.DeleteAPIEndpoint)
		api.GET("/api-endpoints/:apiVersionID", h.GetAllAPIEndpointsByAPIVersionID)
	}

}
