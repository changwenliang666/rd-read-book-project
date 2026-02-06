package service

import (
	"errors"
	"rd-read-book-project/config"
	"rd-read-book-project/model"
)

type userVo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type userInputJson struct {
	Username string `json:"username" binding:"required"`
}

type UserRegisterJson struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(userInfo *UserRegisterJson) error {
	var existCount int64 = 0
	err := config.DB.Model(model.User{}).Where("username=?", userInfo.Username).Count(&existCount).Error
	if err != nil { // 数据库查询失败
		return err
	}
	if existCount > 0 {
		return errors.New("该用户已经存在")
	}

	if err := config.DB.Model(model.User{}).Create(&userInfo).Error; err != nil {
		return errors.New("注册用户失败")
	}

	return nil
}

func GetUserInfoById(userId string) (any, string) {
	var user = userVo{}

	if err := config.DB.Model(model.User{}).Omit("password").First(&user, userId).Error; err != nil {
		errMsg := "未找到该用户信息"
		return nil, errMsg
	}

	return user, ""
}

func UpdateUserName(userInputId int, userJson userInputJson) (string, error) {

	var user model.User
	var responceMsg string = "更新数据成功"

	if err := config.DB.Model(model.User{}).First(&user, userInputId).Error; err != nil {
		responceMsg = "该用户不存在"
		return responceMsg, err
	}

	user.Username = userJson.Username // 更新用户名

	result := config.DB.Model(&model.User{}).Where("id = ?", userInputId).Updates(map[string]interface{}{"username": userJson.Username})

	if result.Error != nil || result.RowsAffected == 0 {
		responceMsg = "更新数据失败"
		return responceMsg, result.Error
	}
	return responceMsg, nil
}
