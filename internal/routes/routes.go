package routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ptm/internal/db/seeder"
	"ptm/internal/middlewares"
	"ptm/internal/utils/response"
)

func InitRoutes(e *echo.Echo) {

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.JWTAuthenticate())
	e.Use(middlewares.ErrorMiddleware())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})
	e.POST("/seeder", func(c echo.Context) error {
		seeder.SeedUsers()
		return response.Ok(c, "OK", nil)
	})
	e.GET("/test", func(c echo.Context) error {
		return errors.New("ERROR SENT")
	})
	RegisterUserRoutes(e)
	RegisterAuthRoutes(e)
}
