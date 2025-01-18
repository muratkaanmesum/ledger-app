package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterAuthRoutes(e *echo.Group) {
	route := e.Group("/auth")

	route.POST("/register", controllers.RegisterUser)
	route.POST("/login", controllers.AuthenticateUser)
}
