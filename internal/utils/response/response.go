package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Ok(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c echo.Context, message string, err error) error {
	return c.JSON(http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Message: message,
		Error:   formatError(err),
	})
}

func NotFound(c echo.Context, message string, err error) error {
	return c.JSON(http.StatusNotFound, Response{
		Status:  http.StatusNotFound,
		Message: message,
		Error:   formatError(err),
	})
}

func InternalServerError(c echo.Context, message string, err error) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Status:  http.StatusInternalServerError,
		Message: message,
		Error:   formatError(err),
	})
}

func Custom(c echo.Context, status int, message string, data any, err error) error {
	return c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   formatError(err),
	})
}

func formatError(err error) any {
	if err == nil {
		return nil
	}
	return err.Error()
}
