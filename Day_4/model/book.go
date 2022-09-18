package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title  string `gorm:"column:title;NOT NULL"`
	Isbn   string `gorm:"column:isbn;NOT NULL"`
	Writer string `gorm:"column:writer;NOT NULL"`
}

func (b *Book) TableName() string {
	return "books"
}

func (b *Book) Create(db *gorm.DB) error {
	err := db.Debug().Create(&b).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *Book) Get(db *gorm.DB) (*Book, error) {
	err := db.Debug().Where(&b).First(&b).Error
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Book) Gets(db *gorm.DB) ([]*Book, error) {
	users := make([]*Book, 0)
	err := db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (b *Book) Update(db *gorm.DB) error {
	err := db.Debug().Model(&b).Where("id = ?", b.ID).Updates(b).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *Book) Delete(db *gorm.DB) error {
	err := db.Debug().Delete(&b, b.ID).Error
	if err != nil {
		return err
	}

	return nil
}

type BookRequest struct {
	Title  string `json:"title" validate:"required,contains,min=4"`
	Isbn   string `json:"isbn" validate:"required,numeric,min=4"`
	Writer string `json:"writer" validate:"required,contains,min=4"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Isbn   string `json:"isbn"`
	Writer string `json:"writer"`
}
