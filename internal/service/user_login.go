package service

import (
	"errors"
	"fmt"
	"rd-read-book-project/config"
	"rd-read-book-project/model"
	"rd-read-book-project/pkg/bcrypt"
	"rd-read-book-project/pkg/jwt"

	"gorm.io/gorm"
)

func UserLogin(username string, password string) (string, error) {
	var user model.User
	res := config.DB.Model(model.User{}).Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("用户名不存在")
		} else {
			return "", fmt.Errorf(res.Error.Error())
		}
	} else {
		if bcrypt.CheckPassword(password, user.Password) { // 密码正确
			token, err := jwt.GenerateToken(username, user.Password)
			if err != nil {
				return "", fmt.Errorf("生成令牌失败")
			} else {
				return token, nil
			}
		} else {
			return "", fmt.Errorf("密码错误")
		}
	}
}
