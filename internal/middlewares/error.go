package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"ptm/internal/utils/customError"
	"ptm/pkg/logger"
)

var errorMessages = map[int]string{
	400: "Bad Request. Please check your input.",
	401: "Unauthorized. Please provide valid credentials.",
	403: "Forbidden. You don't have permission to access this resource.",
	404: "Not Found. The requested resource does not exist.",
	500: "Internal Server Error. Please try again later.",
	503: "Service Unavailable. Please try again later.",
}

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				var customErr *customError.Error
				if errors.As(err, &customErr) {
					message := getErrorMessage(int(customErr.Code))
					logger.Logger.Error("Error Message",
						zap.Int("code", int(customErr.Code)),
						zap.String("message", message),
					)
					return c.JSON(int(customErr.Code), map[string]any{
						"status":  customErr.Code,
						"message": message,
					})
				}

				message := getErrorMessage(500)
				logger.Logger.Error("Generic error encountered", zap.String("message", message))
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"status":  http.StatusInternalServerError,
					"message": message,
				})
			}
			return nil
		}
	}
}

func getErrorMessage(statusCode int) string {
	if message, exists := errorMessages[statusCode]; exists {
		return message
	}
	return "An unexpected error occurred."
}
