package service

import (
	"errors"
	"fmt"
	"rd-read-book-project/config"
	"rd-read-book-project/internal/vo"
	"rd-read-book-project/model"
	"rd-read-book-project/pkg/bcrypt"

	"gorm.io/gorm"
)

type userInputJson struct {
	Username string `json:"username" binding:"required"`
}

type UserRegisterJson struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(userInfo *UserRegisterJson) error {
	var user model.User
	res := config.DB.Model(model.User{}).Where("username=?", userInfo.Username).First(&user)
	if res.Error == nil {
		return fmt.Errorf("用户名已经存在，请重新输入")
	} else if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		encodePassword, encodeError := bcrypt.HashPassword(userInfo.Password)
		if encodeError != nil {
			return encodeError
		}
		newUser := model.User{
			Username: userInfo.Username,
			Password: encodePassword,
		}
		create_error := config.DB.Model(model.User{}).Create(&newUser)
		if create_error.Error != nil || create_error.RowsAffected == 0 {
			return fmt.Errorf("注册用户失败")
		}
		return nil
	} else if res.Error != nil {
		return fmt.Errorf("数据库操作失败")
	}

	return nil
}

func GetUserInfoById(userId string) (any, error) {
	var user model.User
	res := config.DB.Model(model.User{}).Where("id=?", userId).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		errMsg := "未找到该用户信息"
		return nil, fmt.Errorf(errMsg)
	} else if res.Error != nil {
		essMsg := "查询数据库出错"
		return nil, fmt.Errorf(essMsg)
	}

	userInfoVo := vo.UserInfoVo{
		Id:       user.Id,
		Username: user.Username,
	}

	return userInfoVo, nil
}

func UpdateUserName(userInputId int, userJson userInputJson) error {
	var user model.User

	res := config.DB.Model(model.User{}).First(&user, userInputId)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		errorMsg := "用户不存在"
		return fmt.Errorf(errorMsg)
	} else if res.Error != nil {
		errorMsg := "数据库查询失败"
		return fmt.Errorf(errorMsg)
	} else if user.Username == userJson.Username {
		return fmt.Errorf("新名字和当前用户名一致")
	}

	result := config.DB.Model(&model.User{}).Where("id = ?", userInputId).Updates(map[string]interface{}{"username": userJson.Username})

	if result.Error != nil {
		errorMsg := "数据库更新失败"
		return fmt.Errorf(errorMsg)
	} else if result.RowsAffected == 0 {
		errorMsg := "数据更新失败或用户名未发生变化"
		return fmt.Errorf(errorMsg)
	}

	return nil
}
