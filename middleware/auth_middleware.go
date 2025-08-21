package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware authenticates requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Authorization header required"}})
			c.Abort()
			return
		}

		// Bearer Token
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Invalid Authorization header format"}})
			c.Abort()
			return
		}

		tokenString = parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("your_secret_key"), nil // Use the same secret key as in Login
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Invalid or expired token"}})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Invalid token claims"}})
			c.Abort()
			return
		}

		// Set user_id in context
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}
