package controller

import (
	"rd-read-book-project/internal/service"
	"rd-read-book-project/pkg/response"

	"github.com/gin-gonic/gin"
)

func UserLogin(ctx *gin.Context) {
	var user = struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.Fail(ctx, "参数错误", nil)
		return
	}
	token, err := service.UserLogin(user.Username, user.Password)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "登录成功", struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
