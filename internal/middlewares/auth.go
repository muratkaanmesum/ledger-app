package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RoleBasedAuthorization(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			if userRole == nil {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
			}

			if userRole != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "You do not have the required permissions"})
			}

			return next(c)
		}
	}
}
