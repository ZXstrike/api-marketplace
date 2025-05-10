package jwt

import (
	"fmt"
	"time"

	"github.com/ZXstrike/config"
	"github.com/ZXstrike/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	signedToken, err := token.SignedString(config.Config.PrivateKey)

	if err != nil {
		fmt.Println("Error while signing the token: ", err)
		return "", err
	}

	return signedToken, nil
}
