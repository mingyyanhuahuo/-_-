package middleware

import (
	"errors"
	"lostfound/pkg/errcode"
	"lostfound/pkg/jwtutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const perfix = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, perfix))
		if !strings.HasPrefix(header, perfix) || token == "" {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		claims, err := jwtutil.ParseToken(token, jwtutil.TokenTypeAccess)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.Error(errcode.ErrTokenExpired)
			} else {
				c.Error(errcode.ErrUnauthorized)
			}
			c.Abort()
			return
		}
		c.Set("id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("id")
		if !exists {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}
		c.Error(errcode.ErrForbidden)
		c.Abort()
	}
}
