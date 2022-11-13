package http

import (
	"api-mvc/internal/app/auth"
	"api-mvc/internal/app/book"
	"api-mvc/internal/app/user"
	"api-mvc/internal/factory"
	m "api-mvc/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func NewHttp(e *echo.Echo, f *factory.Factory) {
	m.LogMiddleware(e)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))

	api := e.Group("/api")
	auth.NewHandler(f).Route(api.Group("/auth"))
	user.NewHandler(f).Route(api.Group("/users"))
	book.NewHandler(f).Route(api.Group("/books"))

}
