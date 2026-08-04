package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		HeaderValue := c.GetHeader("Authorization")
		if HeaderValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized user",
			})
			c.Abort()
			return
		}

		parts := strings.Split(HeaderValue, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "something went wrong",
			})
			c.Abort()
			return
		}
		tokeString := parts[1]
		token, err := jwt.Parse(tokeString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil

		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid and expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()

			return
		}
		userID := int(claims["user_id"].(float64))

		c.Set("userID", userID)
		c.Next()

	}
}
