package middleware

import (
	"api-mvc/web"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccessControl() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Get("username").(string)
			fmt.Printf("username: %v", username)
			if username != "admin" {
				return c.JSON(http.StatusUnauthorized, web.ResponError{
					Status:  http.StatusUnauthorized,
					Message: errors.New("only admin can access").Error(),
				})
			}

			return next(c)
		}
	}
}
