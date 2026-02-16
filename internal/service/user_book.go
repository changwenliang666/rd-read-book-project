package service

import (
	"errors"
	"fmt"
	"rd-read-book-project/config"
	"rd-read-book-project/model"

	"gorm.io/gorm"
)

func GetBookList(userId int) ([]model.Book, error) {
	var userBookList []model.UserBook
	var userBookId []int
	var bookList []model.Book

	result := config.DB.Model(&model.UserBook{}).Where("user_id = ?", userId).Find(&userBookList)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []model.Book{}, nil
		} else {
			return []model.Book{}, fmt.Errorf("查询用户图书失败")
		}
	}

	for _, bookInfo := range userBookList {
		userBookId = append(userBookId, bookInfo.BookId)
	}

	bookQueryResult := config.DB.Model(&model.Book{}).Where("id IN ?", userBookId).Find(&bookList)

	if bookQueryResult.Error != nil {
		if errors.Is(bookQueryResult.Error, gorm.ErrRecordNotFound) {
			return []model.Book{}, nil
		} else {
			return []model.Book{}, fmt.Errorf("查询用户图书失败")
		}
	}

	return bookList, nil
}
