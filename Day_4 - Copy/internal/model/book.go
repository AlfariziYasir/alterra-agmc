package model

type Book struct {
	Base
	Title  string `gorm:"column:title;NOT NULL"`
	Isbn   string `gorm:"column:isbn;NOT NULL"`
	Writer string `gorm:"column:writer;NOT NULL"`
}

func (b *Book) TableName() string {
	return "books"
}
