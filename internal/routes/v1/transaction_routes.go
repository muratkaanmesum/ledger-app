package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterTransactionRoutes(c *echo.Group) {
	group := c.Group("/transactions")

	group.POST("/credit", controllers.HandleCredit)
	group.POST("/debit", controllers.HandleDebit)
}
