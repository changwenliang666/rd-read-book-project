package controller

import (
	"fmt"
	"rd-read-book-project/internal/service"
	"rd-read-book-project/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user service.UserRegisterJson
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.Fail(ctx, response.ResponseErrorCode.ParamsCode, "参数格式不正确", nil)
		return
	}

	if len(user.Username) < 5 || len(user.Username) > 50 {
		response.Fail(ctx, response.ResponseErrorCode.ParamsCode, "用户名的长度要求在5~50个字符之间", nil)
		return
	}

	if len(user.Password) < 8 || len(user.Password) > 100 {
		response.Fail(ctx, response.ResponseErrorCode.ParamsCode, "密码长度需要介于8 ~ 100 之间", nil)
		return
	}

	err := service.Register(&user)

	if err != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, err.Error(), nil)
		return
	}

	response.Success(ctx, "用户注册成功", nil)

}

func GetUserInfo(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		response.Fail(ctx, response.ResponseErrorCode.VerifCode, "获取身份信息失败", nil)
		return
	}
	result, err := service.GetUserInfoById(userID.(int))
	if err != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, err.Error(), nil)
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
		response.Fail(ctx, response.ResponseErrorCode.ParamsCode, "id无效", nil)
		return
	}

	if err2 := ctx.ShouldBindJSON(&userInputJson); err2 != nil {
		response.Fail(ctx, response.ResponseErrorCode.ParamsCode, "参数错误", nil)
		return
	}
	fmt.Printf("inputInfo --- %d %v", userInputId, userInputJson.Username)
	serviceError := service.UpdateUserName(userInputId, userInputJson)

	if serviceError != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, serviceError.Error(), nil)
		return
	}
	response.Success(ctx, "用户名更新成功", nil)
}
