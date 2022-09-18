package routes

import (
	"api-mvc/controller"
	m "api-mvc/middleware"

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
	m.LogMiddleware(router)
	// recover
	router.Use(middleware.Recover())
	//CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	r := router.Group("/v1")
	r.POST("/register", c.CreateUser)
	r.POST("/login", c.Login)

	user := r.Group("/users", m.SetupAuth())
	user.GET("/refresh", c.Refresh)
	user.GET("/logout", c.Logout)
	user.GET("/:id", c.GetUser)
	user.GET("/", c.GetUsers)
	user.PUT("/:id", c.UpdateUser, m.AccessControl())
	user.DELETE("/:id", c.DeleteUser, m.AccessControl())

	book := r.Group("/books", m.SetupAuth())
	book.POST("/", c.CreateBook)
	book.GET("/:id", c.GetBook)
	book.GET("/", c.GetBooks)
	book.PUT("/:id", c.UpdateBook)
	book.DELETE("/:id", c.DeleteBook)

	return router
}
