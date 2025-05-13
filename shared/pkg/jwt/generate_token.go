package jwt

import (
	"crypto/ecdsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken creates a JWT with custom claims and signs it using the given ECDSA private key.
// It automatically includes an "exp" claim set to 30 days from now if not already provided.
func GenerateAccessToken(privateKey *ecdsa.PrivateKey, customClaims map[string]interface{}) (string, error) {
	// Ensure the "exp" claim is present
	if _, exists := customClaims["exp"]; !exists {
		customClaims["exp"] = time.Now().Add(30 * 24 * time.Hour).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(customClaims))
	return token.SignedString(privateKey)
}
