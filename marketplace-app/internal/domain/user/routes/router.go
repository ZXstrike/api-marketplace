package routes

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/user/handler"
	"github.com/ZXstrike/marketplace-app/internal/domain/user/repositories"
	"github.com/ZXstrike/marketplace-app/internal/domain/user/service"
	"github.com/ZXstrike/marketplace-app/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	repo := repositories.New(db)
	fileRepo := repositories.NewLocalFileRepository("/files")
	service := service.New(repo, fileRepo, privateKey, publicKey)
	h := handler.New(service)

	user := r.Group("/user")
	{
		user.GET("/:id", h.UserProfileHandler)
		user.PUT("/update", middleware.AuthMiddleware(publicKey), h.UpdateUserProfileHandler)
		user.PUT("/change-password", middleware.AuthMiddleware(publicKey), h.ChangePasswordHandler)
		user.POST("/update-profile-picture", middleware.AuthMiddleware(publicKey), h.UpdateProfilePictureHandler)
		user.GET("/profile-picture/:id", h.GetUserProfilePictureHandler)
	}
}
