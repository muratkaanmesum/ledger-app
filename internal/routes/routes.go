package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ptm/internal/db/seeder"
	"ptm/internal/middlewares"
	v1 "ptm/internal/routes/v1"
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

	v1.HandleV1Routes(e)
}
