package router

import (
	"rd-read-book-project/config"
	"rd-read-book-project/internal/logger"
	"rd-read-book-project/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	logger.InitLogger() // 初始化日志

	config.InitDB() // 初始化数据库

	router := gin.Default() // 初始化路由

	logger.Info("初始化路由")
	// 接口使用中间件 cors
	router.Use(middleware.CorsMiddleware())

	InitLoginRouter(router) // login 接口
	InitPingRouter(router)  // ping 接口
	InitUserRouter(router)  // user 接口
	InitBookRouter(router)  // book 接口

	router.Run(":3000")
	return router
}
