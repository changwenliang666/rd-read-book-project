package model

type Book struct {
	ID          int64  `json:"id"`
	Name        string `json:"book_name"`
	Description string `json:"book_description"`
	Author      string `json:"author"`
	Cover       string `json:"cover"`
}

func (Book) TableName() string {
	return "book"
}
