package controller

import (
	"fmt"
	"rd-read-book-project/internal/service"
	"rd-read-book-project/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetBookList(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		response.Fail(ctx, response.ResponseErrorCode.VerifCode, "获取用户信息失败", nil)
		return
	}
	fmt.Println(userID)

	result, err := service.GetBookList(userID.(int))
	if err != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, err.Error(), nil)
		return
	}
	response.Success(ctx, "查询成功", result)
}
