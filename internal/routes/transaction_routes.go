package routes

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterTransactionRoutes(c *echo.Echo) {
	group := c.Group("/transactions")

	group.POST("/credit", controllers.HandleCredit)
}
