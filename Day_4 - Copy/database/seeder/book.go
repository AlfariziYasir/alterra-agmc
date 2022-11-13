package seeder

import (
	"api-mvc/internal/model"
	"log"
	"time"

	"gorm.io/gorm"
)

func bookSeeder(db *gorm.DB) {
	now := time.Now()
	var users = []model.Book{
		{
			Base: model.Base{
				CreatedBy: "admin",
				CreatedAt: now,
				UpdatedBy: "admin",
				UpdatedAt: now,
			},
			Title:  "test book",
			Isbn:   "test@test.com",
			Writer: "test-1",
		},
		{
			Base: model.Base{
				CreatedBy: "admin",
				CreatedAt: now,
				UpdatedBy: "admin",
				UpdatedAt: now,
			},
			Title:  "test book 1",
			Isbn:   "test1@test.com",
			Writer: "test-1",
		},
	}

	err := db.Create(&users).Error
	if err != nil {
		log.Printf("cannot seed data users, with error %v\n", err)
	}
	log.Println("success seed data users")
}
