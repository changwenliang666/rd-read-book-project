package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var successCode = 0
var failCode = -1

func Success(ctx *gin.Context, msg string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    successCode,
		Message: msg,
		Data:    data,
	})
}

func Fail(ctx *gin.Context, msg string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    failCode,
		Message: msg,
		Data:    data,
	})
}
