package router

import (
	"fmt"
	"net/http"
	"rd-read-book-project/config"
	"rd-read-book-project/internal/controller"
	"rd-read-book-project/internal/middleware"
	"rd-read-book-project/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
}

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.JWTAuthMiddleware())
	{
		//根据用户id 查询用户
		userRouter.GET("/get-user-info-pro", controller.GetUserInfo)

		// 根据用户id 更新用户名
		userRouter.PATCH("/updateUserName/:id", controller.UpdateUserName)

		//// 根据用户id 更新用户信息
		//userRouter.PATCH("/updateUserInfo/:id", controller.UpdateUserName)

		// 根据用户id 删除用户
		userRouter.DELETE("/deleteUser/:id", func(ctx *gin.Context) {
			userInputId, err := strconv.Atoi(ctx.Param("id"))
			var user model.User
			if err != nil {
				ctx.JSON(400, gin.H{
					"code":    400,
					"message": "参数错误",
				})
			}
			fmt.Printf("%d", userInputId)
			if err := config.DB.Model(model.User{}).First(&user, userInputId).Error; err != nil {
				ctx.JSON(400, gin.H{
					"code":    400,
					"message": "该用户不存在",
				})
				return
			}

			result := config.DB.Delete(model.User{}, userInputId)
			if result.Error != nil || result.RowsAffected == 0 {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "删除用户失败",
				})
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "删除用户成功",
			})
		})
	}
}
