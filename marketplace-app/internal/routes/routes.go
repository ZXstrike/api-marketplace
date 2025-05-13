package routes

import (
	"github.com/ZXstrike/internal/domain/auth"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	// Auth routes
	auth.AuthRoutes(router)
	// // User routes
	// user := router.Group("/user")
	// {
	// 	user.GET("/:id", GetUserHandler)
	// 	user.PUT("/:id", UpdateUserHandler)
	// 	user.DELETE("/:id", DeleteUserHandler)
	// }

	// // Product routes
	// product := router.Group("/product")
	// {
	// 	product.GET("/", GetProductsHandler)
	// 	product.GET("/:id", GetProductHandler)
	// 	product.POST("/", CreateProductHandler)
	// 	product.PUT("/:id", UpdateProductHandler)
	// 	product.DELETE("/:id", DeleteProductHandler)
	// }
}
