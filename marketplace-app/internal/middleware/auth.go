package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ZXstrike/internal/config"
	"github.com/ZXstrike/shared/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware enforces a valid JWT and makes claims available in context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2) Verify signature & parse claims
		claims, err := jwt.VerifyAccessToken(config.Config.PublicKey, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 3) Check exp claim
		if expVal, ok := claims["exp"].(float64); ok {
			if int64(expVal) < time.Now().Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				return
			}
		}

		// 4) Extract user_id claim
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id claim missing"})
			return
		}

		// 5) Inject into context for handlers:
		//    handlers can retrieve via c.Get("userID")
		c.Set("userID", userID)
		c.Next()
	}
}
