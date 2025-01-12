package middlewares

import (
	"ptm/internal/utils/jwt"
	"ptm/internal/utils/response"

	"github.com/labstack/echo/v4"
)

func RoleBasedAuthorization(requiredRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			claims, err := jwt.ValidateJWT(authHeader)

			if err != nil {
				return response.Forbidden(c, "Forbidden Source", err)
			}
			if authHeader == "" {
				return response.Forbidden(c, "Forbidden Source", nil)
			}

			for _, role := range requiredRoles {
				if claims.Role == role {
					return next(c)
				}
			}

			return response.Forbidden(c, "Not Authorized", nil)
		}
	}
}
