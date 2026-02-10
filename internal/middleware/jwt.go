package middleware

import (
	"rd-read-book-project/pkg/response"
	"strings"

	jwtUtil "rd-read-book-project/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, response.ResponseErrorCode.VerifCode, "缺少 Authorization 头", nil)
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, response.ResponseErrorCode.VerifCode, "Authorization格式错误", nil)
			return
		}

		claims, err := jwtUtil.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, response.ResponseErrorCode.VerifCode, "登录过期", nil)
			return
		}

		// ⭐ 把用户信息放进 Context
		c.Set("user_id", claims.Id)
		c.Set("username", claims.Username)
		c.Next()
	}
}
