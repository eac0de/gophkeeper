package inmiddlewares

import (
	"net/http"
	"strings"

	"auth/internal/services"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(sessionService *services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		clearToken, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		claims, err := sessionService.ParseToken(clearToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set(gin.AuthUserKey, claims.UserID)
		c.Next()
	}
}
