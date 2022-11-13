package seeder

import (
	"api-mvc/internal/model"
	"log"
	"time"

	"gorm.io/gorm"
)

func userSeeder(db *gorm.DB) {
	now := time.Now()
	var users = []model.User{
		{
			Base: model.Base{
				CreatedBy: "admin",
				CreatedAt: now,
				UpdatedBy: "admin",
				UpdatedAt: now,
			},
			Name:     "test",
			Email:    "test@test.com",
			Password: "12345678",
		},
		{
			Base: model.Base{
				CreatedBy: "admin",
				CreatedAt: now,
				UpdatedBy: "admin",
				UpdatedAt: now,
			},
			Name:     "test1",
			Email:    "test1@test.com",
			Password: "12345678",
		},
	}

	err := db.Create(&users).Error
	if err != nil {
		log.Printf("cannot seed data users, with error %v\n", err)
	}
	log.Println("success seed data users")
}
