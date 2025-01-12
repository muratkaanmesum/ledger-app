package routes

import (
	"github.com/labstack/echo/v4"
	"ptm/controllers"
)

func RegisterUserRoutes(e *echo.Echo) {
	userRoute := e.Group("/users")

	userRoute.GET("/:id", controllers.GetUserById)
	userRoute.GET("/", controllers.GetAllUsers)
}
