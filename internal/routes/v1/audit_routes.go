package v1

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/controllers"
	"ptm/internal/middlewares"
)

func RegisterAuditRoutes(e *echo.Group) {
	route := e.Group("/audit")

	route.Use(middlewares.RoleBasedAuthorization("admin"))

	controller := controllers.NewAuditLogController()
	route.POST("", controller.GetLogs)
	route.GET("/:id", controller.GetLogsById)
}
