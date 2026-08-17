package middlewares

import (
	"strings"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/service"

	"github.com/gin-gonic/gin"
)

type AuthenticateMiddleware struct {
	AS service.AuthService
}

func (AM *AuthenticateMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") == "" {
			c.Error(apperrors.ErrUnauthorized)
			c.Abort()
			return
		}
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.Error(apperrors.ErrNotSupportedAuthenticateMethod)
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := AM.AS.VerifyToken(tokenString)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("Role", claims.Role)
		c.Next()
	}
}
