package middlewares

import (
	"ptm/internal/utils/response"
	"ptm/pkg/jwt"
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

func JWTAuthenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			publicRoutes := map[string]bool{
				"/health":               true,
				"/api/v1/auth/login":    true,
				"/api/v1/auth/register": true,
			}

			if publicRoutes[c.Path()] {
				return next(c)
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Forbidden(c, "Token doesn't exist", nil)
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return response.Forbidden(c, "Token is invalid", nil)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwt.ValidateJWT(tokenString)
			if err != nil {
				return response.Forbidden(c, "Token is invalid", nil)
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}
