package middleware

import (
	"strings"

	"github.com/alexe0110/chat-system/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, "not token in headers")
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := auth.ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(401, "not auth")
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
