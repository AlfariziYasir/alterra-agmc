package controller

import (
	"api-mvc/db"

	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func NewController() *Controller {
	db, _ := db.NewClient()

	return &Controller{db.Conn()}
}
