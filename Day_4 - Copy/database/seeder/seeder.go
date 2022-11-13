package seeder

import (
	"api-mvc/database/postgres"

	"gorm.io/gorm"
)

type seed struct {
	DB *gorm.DB
}

func NewSeeder() *seed {
	db, _ := postgres.NewClient()
	return &seed{db.Conn()}
}

func (s *seed) SeedAll() {
	userSeeder(s.DB)
	bookSeeder(s.DB)
}

func (s *seed) DeleteAll() {
	s.DB.Exec("DELETE FROM users")
	s.DB.Exec("DELETE FROM books")
}
