package router

import (
	"rd-read-book-project/internal/controller"
	"rd-read-book-project/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitBookRouter(router *gin.Engine) {
	bookRouter := router.Group("/book")
	bookRouter.Use(middleware.JWTAuthMiddleware())
	{
		bookRouter.GET("get-book-list", controller.GetBookList)
		bookRouter.POST("create-book", controller.CreateBook)
	}

}
