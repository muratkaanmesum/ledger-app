package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"ptm/internal/utils/logger"
)

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logger.Logger.Error("Unhandled error in middleware", zap.Error(err))

				var httpError *echo.HTTPError
				if errors.As(err, &httpError) {
					return c.JSON(httpError.Code, map[string]interface{}{
						"status":  "error",
						"message": httpError.Message,
					})
				}

				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"status":  "error",
					"message": "Internal server error",
				})
			}

			return nil
		}
	}
}
