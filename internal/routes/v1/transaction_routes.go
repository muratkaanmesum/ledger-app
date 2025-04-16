package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
	"ptm/pkg/worker"
)

func RegisterTransactionRoutes(c *echo.Group) {
	controller := controllers.NewTransactionController()
	group := c.Group("/transactions")
	poolName := "transactions"

	group.POST("/history", controller.GetTransactions)
	group.POST("/credit", worker.RunWithWorker[controllers.CreditRequest](controller.HandleCredit, poolName))
	group.POST("/debit", worker.RunWithWorker[controllers.CreditRequest](controller.HandleDebit, poolName))
	group.POST("/transfer", worker.RunWithWorker[controllers.TransferRequest](controller.HandleTransfer, poolName))
	group.POST("/schedule", worker.RunWithWorker[controllers.ScheduleRequest](controller.ScheduleTransaction, poolName))
	group.POST("/:id", controller.GetById)

}
