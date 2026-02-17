package service

import (
	"errors"
	"fmt"
	"rd-read-book-project/config"
	"rd-read-book-project/model"
	"rd-read-book-project/pkg/epub"

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

func CreateBook(meta *epub.Metadata, userId int) error {
	var book model.Book
	var userBook model.UserBook

	book.Author = meta.Author
	book.Name = meta.Title
	book.Description = meta.Description
	book.Cover = ""
	book.RemoteUrl = meta.RemoteUrl

	return config.DB.Transaction(func(tx *gorm.DB) error {
		bookResult := tx.Model(&model.Book{}).Create(&book)

		if bookResult.Error != nil || bookResult.RowsAffected == 0 {
			return fmt.Errorf("创建书籍失败")
		}

		userBook.BookId = book.ID
		userBook.UserId = userId

		userBookResult := tx.Model(&model.UserBook{}).Create(&userBook)

		if userBookResult.Error != nil || userBookResult.RowsAffected == 0 {
			return fmt.Errorf("创建书籍失败")
		}

		return nil
	})
}
