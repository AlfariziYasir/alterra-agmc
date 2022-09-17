package middleware

import (
	"api-mvc/db/redis"
	"api-mvc/lib/token"
	"api-mvc/web"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r, _ := redis.NewClient()

			data, err := token.TokenValid(c.Request())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, web.ResponError{
					Status:  http.StatusUnauthorized,
					Message: errors.New("not allowed to access").Error(),
				})
			}

			_, err = r.Conn().Get(context.Background(), data.TokenUuid).Result()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, web.ResponError{
					Status:  http.StatusInternalServerError,
					Message: "token invalid",
				})
			}

			c.Set("user_id", data.UserId)
			c.Set("username", data.Username)

			return next(c)
		}
	}
}
