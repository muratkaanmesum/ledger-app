package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"ptm/pkg/logger"
	"ptm/pkg/utils/customError"
)

var errorMessages = map[int]string{
	http.StatusBadRequest:          "Invalid request. Please check your input.",
	http.StatusUnauthorized:        "Authentication required.",
	http.StatusForbidden:           "Access forbidden.",
	http.StatusNotFound:            "Resource not found.",
	http.StatusInternalServerError: "An unexpected error occurred.",
	http.StatusServiceUnavailable:  "Service unavailable. Please try again later.",
}

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				var customErr *customError.Error
				if ok := errors.As(err, &customErr); ok {
					return handleError(c, customErr)
				}

				logger.Logger.Error("Unhandled error",
					zap.String("error", err.Error()),
				)

				return err
			}
			return nil
		}
	}
}

func handleError(c echo.Context, customErr *customError.Error) error {
	logger.Logger.Error("Custom error occurred",
		zap.Int("code", int(customErr.Code)),
		zap.String("message", customErr.Message),
		zap.Error(customErr.Details),
	)

	return c.JSON(int(customErr.Code), map[string]interface{}{
		"status":  customErr.Code,
		"message": customErr.Message,
		"details": customErr.Details,
	})
}

func getErrorMessage(statusCode int) string {
	if message, exists := errorMessages[statusCode]; exists {
		return message
	}
	return "An unexpected error occurred."
}
