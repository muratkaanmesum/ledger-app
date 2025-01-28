package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ptm/internal/db/seeder"
	"ptm/internal/middlewares"
	v1 "ptm/internal/routes/v1"
	"ptm/internal/utils/response"
)

func InitRoutes(e *echo.Echo) {

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middlewares.JWTAuthenticate())
	e.Use(middlewares.ErrorMiddleware())
	e.Use(middlewares.PerformanceMiddleware())

	e.POST("/seeder", func(c echo.Context) error {
		seeder.SeedUsers()
		return response.Ok(c, "OK")
	})

	v1.HandleV1Routes(e)
}
