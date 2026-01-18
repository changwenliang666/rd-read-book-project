package router

import (
	"encoding/json"
	"net/http"

	jwtUtil "rd-read-book-project/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
}

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		// 注册
		userRouter.POST("/register", func(ctx *gin.Context) {
			body, _ := ctx.GetRawData()
			var user = map[string]any{}
			json.Unmarshal(body, &user)
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "注册成功",
			})
		})
		// 登录
		userRouter.POST("/login", func(ctx *gin.Context) {
			var user struct {
				Username string `json:"username" binding:"required"`
				UserId   uint   `json:"userId" binding:"required"`
			}

			if err := ctx.ShouldBindJSON(&user); err != nil {
				ctx.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}

			token, err := jwtUtil.GenerateToken(user.UserId, user.Username)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "用户信息认证失败",
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code":  200,
				"token": token,
			})
		})
	}
}
