package factory

import (
	"api-mvc/database/postgres"
	"api-mvc/database/redis"
	"api-mvc/internal/repository"
	l "api-mvc/pkg/logger"
)

type Factory struct {
	Auth repository.Auth
	User repository.User
	Book repository.Book
}

func NewFactory() *Factory {
	pg, err := postgres.NewClient()
	if err != nil {
		l.Log().Err(err).Msg(err.Error())
	}

	rds, err := redis.NewClient()
	if err != nil {
		l.Log().Err(err).Msg(err.Error())
	}

	return &Factory{
		repository.NewAuthRepository(rds.Conn()),
		repository.NewUserRepository(pg.Conn(), rds.Conn()),
		repository.NewBookRepository(pg.Conn()),
	}
}
