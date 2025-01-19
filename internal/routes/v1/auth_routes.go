package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterAuthRoutes(e *echo.Group) {
	route := e.Group("/auth")

	controller := controllers.NewAuthController()

	route.POST("/register", controller.RegisterUser)
	route.POST("/login", controller.AuthenticateUser)
}
