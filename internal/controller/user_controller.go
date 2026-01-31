package controller

import (
	"rd-read-book-project/internal/service"
	"rd-read-book-project/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
