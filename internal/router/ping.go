package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

func InitPingRouter(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
			"time":    time.Now().UnixMilli(),
		})
	})
}
