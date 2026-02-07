package router

import (
	"rd-read-book-project/internal/controller"

	"github.com/gin-gonic/gin"
)

func InitLoginRouter(router *gin.Engine) {
	router.POST("/login", controller.UserLogin)
	router.POST("/register", controller.Register)
}
