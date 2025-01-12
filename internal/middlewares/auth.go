package middlewares

import (
	"net/http"
	"ptm/internal/utils/jwt"
	"ptm/internal/utils/response"
	"strings"

	"github.com/labstack/echo/v4"
)

func RoleBasedAuthorization(requiredRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(*jwt.CustomClaims)
			if !ok {
				return response.Forbidden(c, "User data not found in context", nil)
			}

			for _, role := range requiredRoles {
				if userClaims.Role == role {
					return next(c)
				}
			}

			return response.Forbidden(c, "You are not authorized to access this resource", nil)
		}
	}
}

func JWTAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.Forbidden(c, "Token doesnt exist", nil)
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Forbidden(c, "Token is invalid", nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
		}

		c.Set("user", claims)

		return next(c)
	}
}
