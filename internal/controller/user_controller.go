package controller

import (
	"rd-read-book-project/internal/service"
	"rd-read-book-project/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user service.UserRegisterJson
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.Fail(ctx, "参数格式不正确", nil)
		return
	}

	if len(user.Username) < 5 || len(user.Username) > 50 {
		response.Fail(ctx, "用户名的长度要求在5~50个字符之间", nil)
		return
	}

	if len(user.Password) < 8 || len(user.Password) > 100 {
		response.Fail(ctx, "用户名过长", nil)
		return
	}

	err := service.Register(&user)

	if err != nil {
		response.Fail(ctx, "用户注册失败", nil)
		return
	}

	response.Success(ctx, "用户注册成功", nil)

}

func GetUserInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	result, errMsg := service.GetUserInfoById(id)
	if errMsg != "" {
		response.Fail(ctx, errMsg, nil)
		return
	}
	response.Success(ctx, "查询成功", result)
}

func UpdateUserName(ctx *gin.Context) {
	var userInputJson = struct {
		Username string `json:"username" binding:"required"`
	}{}

	var userIdStr = ctx.Param("id")
	userInputId, err := strconv.Atoi(userIdStr)

	if err != nil {
		response.Fail(ctx, "id无效", nil)
		return
	}
	if err := ctx.ShouldBindJSON(&userInputJson); err != nil {
		response.Fail(ctx, "参数错误", nil)
		return
	}

	result, err := service.UpdateUserName(userInputId, userInputJson)

	if err != nil {
		response.Fail(ctx, result, nil)
		return
	}
	response.Success(ctx, result, nil)

}
