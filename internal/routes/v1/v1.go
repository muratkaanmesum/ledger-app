package v1

import (
	"github.com/labstack/echo/v4"
)

func HandleV1Routes(c *echo.Echo) {

	v1Group := c.Group("/api/v1")

	RegisterUserRoutes(v1Group)
	RegisterAuthRoutes(v1Group)
	RegisterTransactionRoutes(v1Group)
	RegisterBalanceRoutes(v1Group)
	RegisterAuditRoutes(v1Group)
}
