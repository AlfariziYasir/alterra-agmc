package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"column:name;NOT NULL"`
	Email     string `gorm:"column:email;NOT NULL"`
	Password  string `gorm:"column:password;NOT NULL"`
	TokenUuid string `gorm:"column:token_uuid"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(password)

	return nil
}

func (u *User) Create(db *gorm.DB) error {
	temp := new(User)
	err := db.Debug().Session(&gorm.Session{SkipHooks: false}).Create(&u).Scan(&temp).Error
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
	Name      string `json:"name" validate:"required,alphanum,min=4"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,contains,min=6"`
	TokenUuid string
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
