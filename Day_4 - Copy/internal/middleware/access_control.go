package middleware

import (
	"api-mvc/pkg/util/response"
	"errors"

	"github.com/labstack/echo/v4"
)

func AccessControl() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Get("username").(string)

			if username != "admin" {
				return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, errors.New("unauthorized")).Send(c)
			}

			return next(c)
		}
	}
}
