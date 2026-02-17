package router

import (
	"rd-read-book-project/internal/controller"
	"rd-read-book-project/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitUploadRouter(router *gin.Engine) {
	bookRouter := router.Group("/upload")
	bookRouter.Use(middleware.JWTAuthMiddleware())
	{
		bookRouter.POST("upload-file", controller.CreateBook)
	}
}
