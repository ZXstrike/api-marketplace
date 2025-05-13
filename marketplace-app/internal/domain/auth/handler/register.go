package handler

import (
	"github.com/ZXstrike/internal/database"
	"github.com/ZXstrike/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(context *gin.Context) {

	requestBody := struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Check if user already exists

	var user models.User

	if err := database.PostgresDB.Where("email = ?", requestBody.Email).First(&user).Error; err == nil {
		context.JSON(400, gin.H{
			"error": "User already exists",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)

	if err != nil {
		context.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	var newUser models.User

	newUser.Email = requestBody.Email
	newUser.Username = requestBody.Username
	newUser.PasswordHash = string(hash)

	if err := database.PostgresDB.Create(&newUser).Error; err != nil {
		context.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Generate JWT token
	// token, err := generateToken(newUser)
	// if err != nil {
	// 	context.JSON(500, gin.H{
	// 		"error": "Internal server error",
	// 	})
	// 	return
	// }

	// Return the token to the client
	context.JSON(200, gin.H{
		"message": "Registration successful",
		// "token":   token,
	})
}
