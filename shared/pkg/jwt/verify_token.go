package jwt

import (
	"crypto/ecdsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyAccessToken verifies a JWT using the provided ECDSA public key and returns the claims if valid.
func VerifyAccessToken(publicKey *ecdsa.PublicKey, tokenString string) (map[string]interface{}, error) {
	// Parse the token with the expected signing method and key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and assert the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert MapClaims to map[string]interface{}
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, errors.New("invalid token")
}
