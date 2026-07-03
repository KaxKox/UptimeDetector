package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"gocheck/internal/auth"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak tokenu"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Nieprawidłowy token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}