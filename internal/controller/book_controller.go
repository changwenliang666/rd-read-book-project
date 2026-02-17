package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"rd-read-book-project/internal/service"
	epubparser "rd-read-book-project/pkg/epub"
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

func CreateBook(ctx *gin.Context) {
	stfpRes, err := UploadFileSFTP(ctx)
	if err != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, err.Error(), nil)
		return
	}
	r, err := zip.NewReader(bytes.NewReader(stfpRes.Data), int64(len(stfpRes.Data)))
	if err != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, "解析 EPUB 文件失败", nil)
		return
	}
	meta, err := epubparser.ParseEPUBFromZipReader(r)
	if err != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, "解析 EPUB 元信息失败", nil)
		return
	}
	meta.RemoteUrl = stfpRes.RemoteUrl

	userId, ok := ctx.Get("user_id")
	if ok == false {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, "获取用户信息失败", nil)
		return
	}
	serviceErr := service.CreateBook(meta, userId.(int))
	if serviceErr != nil {
		response.Fail(ctx, response.ResponseErrorCode.BaseCode, serviceErr.Error(), nil)
		return
	}
	response.Success(ctx, "添加图书成功", nil)

}
