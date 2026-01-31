package router

import (
	"fmt"
	"net/http"
	"rd-read-book-project/config"
	"rd-read-book-project/model"
	"strconv"

	jwtUtil "rd-read-book-project/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
}

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		// 注册
		userRouter.POST("/register", func(ctx *gin.Context) {
			var user model.User
			if err := ctx.ShouldBindJSON(&user); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "参数传递错误",
				})
				return
			}
			if err := config.DB.Create(&user).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "注册用户失败",
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "success",
			})

		})
		// 登录
		userRouter.POST("/login", func(ctx *gin.Context) {
			var user struct {
				Username string `json:"username" binding:"required"`
				UserId   uint   `json:"userId" binding:"required"`
			}

			if err := ctx.ShouldBindJSON(&user); err != nil {
				ctx.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}

			token, err := jwtUtil.GenerateToken(user.UserId, user.Username)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "用户信息认证失败",
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code":  200,
				"token": token,
			})
		})

		//根据用户id 查询用户
		userRouter.GET("/getUserInfo", func(ctx *gin.Context) {
			id := ctx.Query("id")
			var user = struct {
				Id       int    `json:"id"`
				Username string `json:"username"`
			}{}
			if err := config.DB.Model(model.User{}).Omit("password").First(&user, id).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "未找到该用户信息",
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "success",
				"data": user,
			})
		})

		// 根据用户id 更新用户名
		userRouter.PATCH("/updateUserName/:id", func(ctx *gin.Context) {
			var userIdStr = ctx.Param("id")
			userInputId, err := strconv.Atoi(userIdStr)
			if err != nil {
				ctx.JSON(400, gin.H{"code": 400, "message": "ID 无效"})
				return
			}
			var userInputJson = struct {
				Username string `json:"username" binding:"required"`
			}{}

			var user model.User
			if err := ctx.ShouldBindJSON(&userInputJson); err != nil {
				ctx.JSON(200, gin.H{
					"code":    400,
					"message": "参数名称错误",
				})
				return
			}

			if err := config.DB.Model(model.User{}).First(&user, userInputId).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    400,
					"message": "该用户不存在",
				})
				return
			}

			user.Username = userInputJson.Username // 更新用户名
			// 避坑写法：Save 会更新所有字段
			//if err := config.DB.Model(&user).Save(&user).Error; err != nil {
			//	ctx.JSON(http.StatusOK, gin.H{
			//		"code":    500,
			//		"message": "更新数据失败",
			//	})
			//	return
			//}

			result := config.DB.Model(&model.User{}).Where("id = ?", userInputId).Update("username", userInputJson.Username)
			if result.Error != nil || result.RowsAffected == 0 {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "更新数据失败",
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "更新成功",
			})
		})

		// 根据用户id 更新用户信息
		userRouter.PATCH("/updateUserInfo/:id", func(ctx *gin.Context) {
			var userIdStr = ctx.Param("id")
			userInputId, err := strconv.Atoi(userIdStr)
			if err != nil {
				ctx.JSON(400, gin.H{"code": 400, "message": "ID 无效"})
				return
			}
			var userInputJson = struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}{}

			var user model.User
			if err := ctx.ShouldBindJSON(&userInputJson); err != nil {
				ctx.JSON(200, gin.H{
					"code":    400,
					"message": "参数名称错误",
				})
				return
			}

			if err := config.DB.Model(model.User{}).First(&user, userInputId).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    400,
					"message": "该用户不存在",
				})
				return
			}

			user.Username = userInputJson.Username // 更新用户名
			// 避坑写法：Save 会更新所有字段
			//if err := config.DB.Model(&user).Save(&user).Error; err != nil {
			//	ctx.JSON(http.StatusOK, gin.H{
			//		"code":    500,
			//		"message": "更新数据失败",
			//	})
			//	return
			//}

			result := config.DB.Model(&model.User{}).Where("id = ?", userInputId).Updates(map[string]interface{}{"username": userInputJson.Username, "password": userInputJson.Password})
			if result.Error != nil || result.RowsAffected == 0 {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "更新数据失败",
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "更新成功",
			})
		})

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
