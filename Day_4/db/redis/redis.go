package redis

import (
	"api-mvc/config"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Client interface {
	Conn() *redis.Client
	Close() error
}

func NewClientContext(ctx context.Context) (Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Cfg().RedisHost, config.Cfg().RedisPort),
		DB:   config.Cfg().RedisDB,
	})
	err := db.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &client{db}, nil
}

func NewClient() (Client, error) {
	return NewClientContext(context.Background())
}

type client struct {
	db *redis.Client
}

func (c *client) Conn() *redis.Client { return c.db }
func (c *client) Close() error        { return c.db.Close() }
