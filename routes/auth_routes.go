package routes

import (
	"github.com/labstack/echo/v4"
	"ptm/controllers"
)

func RegisterAuthRoutes(e *echo.Echo) {
	route := e.Group("/auth")

	route.POST("/register", controllers.RegisterUser)
}
