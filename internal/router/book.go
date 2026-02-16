package router

import (
	"rd-read-book-project/internal/controller"
	"rd-read-book-project/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitBookRouter(router *gin.Engine) {
	router.Use(middleware.JWTAuthMiddleware())
	router.GET("get-book-list", controller.GetBookList)
}
