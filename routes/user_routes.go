package routes

import (
	"github.com/labstack/echo/v4"
	"ptm/controllers"
)

func RegisterUserRoutes(e *echo.Echo) {
	userRoute := e.Group("/user")

	userRoute.GET(":userid", controllers.GetUserById)
}
