package model

type UserBook struct {
	Id      int `json:"id"`
	UserId  int `json:"user_id"`
	BookId  int `json:"book_id"`
	Process int `json:"process"`
}

func (user *UserBook) TableName() string {
	return "user_book"
}
