package book

import (
	"api-mvc/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.Use(middleware.SetupAuth())
	g.POST("", h.Create)
	g.GET("", h.Find)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
