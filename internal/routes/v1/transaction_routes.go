package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterTransactionRoutes(c *echo.Group) {
	controller := controllers.NewTransactionController()
	group := c.Group("/transactions")

	group.POST("/history", controller.GetTransactions)
	group.POST("/credit", controller.HandleCredit)
	group.POST("/debit", controller.HandleDebit)
	group.POST("/transfer", controller.HandleTransfer)
	group.POST("/:id", controller.HandleTransfer)
}
