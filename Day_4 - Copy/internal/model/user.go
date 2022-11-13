package model

type User struct {
	Base
	Name     string `gorm:"column:name;NOT NULL"`
	Email    string `gorm:"column:email;NOT NULL"`
	Password string `gorm:"column:password;NOT NULL"`
}

func (u *User) TableName() string {
	return "users"
}
