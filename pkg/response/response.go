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

type ErrorCode struct {
	ParamsCode int // 参数不合法
	VerifCode  int // 身份校验失败
	BaseCode   int // 其他
}

var ResponseErrorCode = ErrorCode{
	ParamsCode: 400,
	VerifCode:  401,
	BaseCode:   -100,
}

func Success(ctx *gin.Context, msg string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    successCode,
		Message: msg,
		Data:    data,
	})
}

func Fail(ctx *gin.Context, code int, msg string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
	ctx.Abort()
}
