package db

import (
	"api-mvc/config"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client interface {
	Conn() *gorm.DB
	Close() error
}

func NewClientContext(ctx context.Context) (Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.Cfg().DBHost,
		config.Cfg().DBPort,
		config.Cfg().DBUser,
		config.Cfg().DBName,
		config.Cfg().DBPass,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})
	if err != nil {
		panic(err)
	}

	test, _ := db.DB()
	err = test.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
	// &model.User{},
	// &model.Book{},
	)
	if err != nil {
		return nil, err
	}

	return &client{db}, nil
}

func NewClient() (Client, error) {
	return NewClientContext(context.Background())
}

type client struct {
	db *gorm.DB
}

func (c *client) Conn() *gorm.DB { return c.db }
func (c *client) Close() error {
	db, _ := c.db.DB()
	return db.Close()
}
