package controller

import (
	redis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Controller struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewController(redis *redis.Client, db *gorm.DB) Controller {
	return Controller{db, redis}
}
