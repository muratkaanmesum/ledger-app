package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Example route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})
}