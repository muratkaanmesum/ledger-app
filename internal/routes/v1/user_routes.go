package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
	"ptm/internal/middlewares"
)

func RegisterUserRoutes(e *echo.Group) {
	c := controllers.NewUserController()

	route := e.Group("/users")

	route.Use(middlewares.RoleBasedAuthorization("admin"))

	route.GET("/:id", c.GetUserById)
	route.PUT("/:id", c.UpdateUser)
	route.DELETE("/:id", c.DeleteUser)
	route.GET("/", c.GetAllUsers)
}
