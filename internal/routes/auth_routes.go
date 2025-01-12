package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ptm/controllers"
	"ptm/internal/middlewares"
)

func RegisterAuthRoutes(e *echo.Echo) {
	route := e.Group("/auth")

	route.POST("/register", controllers.RegisterUser)
	route.POST("/login", controllers.AuthenticateUser)
	route.POST("/test", func(c echo.Context) error {
		fmt.Println("entered")
		return c.String(http.StatusOK, "ok")
	}, middlewares.RoleBasedAuthorization([]string{"admin"}))
}
