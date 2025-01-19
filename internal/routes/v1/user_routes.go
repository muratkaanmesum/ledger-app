package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterUserRoutes(e *echo.Group) {
	c := controllers.NewUserController()

	userRoute := e.Group("/users")

	userRoute.GET("/:id", c.GetUserById)
	userRoute.GET("/", c.GetAllUsers)
}
