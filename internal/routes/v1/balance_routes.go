package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
)

func RegisterBalanceRoutes(c *echo.Group) {
	controller := controllers.NewBalanceController()

	balanceGroup := c.Group("/balances")

	balanceGroup.GET("/current", controller.GetBalance)
	balanceGroup.POST("/at-time", controller.BalanceAtTime)
	balanceGroup.GET("/historical", controller.GetHistoricalBalance)
}
