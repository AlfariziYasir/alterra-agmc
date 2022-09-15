package routes

import (
	"api-mvc/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(c *controller.Controller) *echo.Echo {
	router := echo.New()

	router.GET("/ping", func(c echo.Context) error {
		h := c.Get("User-Agent")

		return c.JSON(200, map[string]interface{}{
			"message": h,
		})
	})

	// logger
	router.Use(middleware.Logger())
	// recover
	router.Use(middleware.Recover())
	//CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	r := router.Group("/v1")
	r.POST("/users", c.CreateUser)
	r.GET("/users/:id", c.GetUser)
	r.GET("/users", c.GetUsers)
	r.PUT("/users/:id", c.UpdateUser)
	r.DELETE("/users/:id", c.DeleteUser)

	r.POST("/books", c.CreateBook)
	r.GET("/books/:id", c.GetBook)
	r.GET("/books", c.GetBooks)
	r.PUT("/books/:id", c.UpdateBook)
	r.DELETE("/books/:id", c.DeleteBook)

	return router
}
