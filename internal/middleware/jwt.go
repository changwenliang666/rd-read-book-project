package middleware

import (
	"net/http"
	"strings"

	jwtUtil "rd-read-book-project/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "缺少 Authorization 头",
			})
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Authorization 格式错误",
			})
			return
		}

		claims, err := jwtUtil.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Token 无效或已过期",
			})
			return
		}

		// ⭐ 把用户信息放进 Context
		c.Set("user_id", claims.Password)
		c.Set("username", claims.Username)
		c.Next()
	}
}
