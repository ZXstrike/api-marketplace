package handler

import (
	"github.com/ZXstrike/internal/database"
	"github.com/ZXstrike/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(context *gin.Context) {

	requestBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	var user models.User

	database.PostgresDB.Where("email = ?", requestBody.Email).First(&user)

	// Check if user exists
	// and if the password is correct

	if user.ID == "" {
		context.JSON(401, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(requestBody.Password)); err != nil {
		context.JSON(401, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	// token, err := generateToken(user)
	// if err != nil {
	// 	context.JSON(500, gin.H{
	// 		"error": "Internal server error",
	// 	})
	// 	return
	// }

	// Return the token to the client
	context.JSON(200, gin.H{
		"message": "Login successful",
		// "token":   token,

	})

}
