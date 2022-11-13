package middleware

import (
	"api-mvc/database/redis"
	"api-mvc/internal/pkg/token"
	"api-mvc/pkg/util/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r, _ := redis.NewClient()

			data, err := token.TokenValid(c.Request())
			if err != nil {
				return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
			}

			_, err = r.Conn().Get(data.TokenUuid).Result()
			if err != nil {
				return response.ErrorBuilder(response.CustomErrorBuilder(http.StatusUnauthorized, "token invalid", "token invalid"), err).Send(c)
			}

			c.Set("user_id", data.UserId)
			c.Set("username", data.Username)
			c.Set("email", data.Email)

			return next(c)
		}
	}
}
