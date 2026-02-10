package service

import (
	"errors"
	"fmt"
	"rd-read-book-project/config"
	"rd-read-book-project/internal/vo"
	"rd-read-book-project/model"
	"rd-read-book-project/pkg/bcrypt"
	"rd-read-book-project/pkg/jwt"

	"gorm.io/gorm"
)

func UserLogin(username string, password string) (vo.LoginUserInfoVo, error) {
	var user model.User
	res := config.DB.Model(model.User{}).Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return vo.LoginUserInfoVo{}, fmt.Errorf("用户名不存在")
		} else {
			return vo.LoginUserInfoVo{}, fmt.Errorf(res.Error.Error())
		}
	} else {
		if bcrypt.CheckPassword(password, user.Password) { // 密码正确
			token, err := jwt.GenerateToken(user.Id, user.Username)
			if err != nil {
				return vo.LoginUserInfoVo{}, fmt.Errorf("生成令牌失败")
			} else {
				return vo.LoginUserInfoVo{
					Token: token,
				}, nil
			}
		} else {
			return vo.LoginUserInfoVo{}, fmt.Errorf("密码错误")
		}
	}
}
