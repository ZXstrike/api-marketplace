package middleware

import (
	"crypto/ecdsa"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(publicKey *ecdsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) < 8 || tokenStr[:7] != "Bearer " {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}
		tokenStr = tokenStr[7:]

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		userID, _ := strconv.Atoi(claims.Subject)
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
