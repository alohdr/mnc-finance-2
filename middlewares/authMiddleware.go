package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mnc-finance/utils"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "UNAUTHORIZED", "message": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check the format of the token (should start with "Bearer")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "UNAUTHORIZED", "message": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify the token
		token, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "UNAUTHORIZED", "message": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user information from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "UNAUTHORIZED", "message": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the request context
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
