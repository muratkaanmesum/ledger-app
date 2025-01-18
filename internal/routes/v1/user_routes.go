package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterUserRoutes(e *echo.Group) {
	userRoute := e.Group("/users")

	userRoute.GET("/:id", controllers.GetUserById)
	userRoute.GET("/", controllers.GetAllUsers)
}
