package jwt

import (
	"fmt"

	"github.com/ZXstrike/config"
	"github.com/golang-jwt/jwt/v5"
)

type ExtractedToken struct {
	UserID uint `json:"user_id"`
}

// VerifyToken function
func VerifyToken(tokenString string) (*ExtractedToken, error) {

	type CustomClaims struct {
		UserID uint `json:"user_id"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Config.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		extractedToken := &ExtractedToken{
			UserID: claims.UserID,
		}

		return extractedToken, nil
	}

	return nil, fmt.Errorf("error while extracting token")
}
