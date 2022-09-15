package model

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"column:name;NOT NULL"`
	Password string `gorm:"column:password;NOT NULL"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Create(db *gorm.DB) error {
	temp := new(User)
	err := db.Debug().Create(&u).Scan(&temp).Error
	fmt.Println("err:", err)
	if err != nil {
		return err
	}

	*u = *temp
	return nil
}

func (u *User) Get(db *gorm.DB) (*User, error) {
	err := db.Debug().Where(&u).First(&u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Gets(db *gorm.DB) ([]*User, error) {
	users := make([]*User, 0)
	err := db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Update(db *gorm.DB) error {
	err := db.Debug().Model(&u).Where("id = ?", u.ID).Updates(u).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete(db *gorm.DB) error {
	err := db.Debug().Delete(&u, u.ID).Error
	if err != nil {
		return err
	}

	return nil
}

type UserRequest struct {
	Name     string `json:"name" validate:"required,alphanum,min=4"`
	Password string `json:"password" validate:"required,contains,min=6"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
