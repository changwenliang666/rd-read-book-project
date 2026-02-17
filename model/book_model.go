package model

type Book struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Cover       string `json:"cover"`
	RemoteUrl   string `json:"remote_url"`
}

func (Book) TableName() string {
	return "book"
}
