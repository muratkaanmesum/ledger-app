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

func Ok(c echo.Context, message string, data ...any) error {
	var responseData any
	if len(data) > 0 {
		responseData = data[0]
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    responseData,
	})
}

func BadRequest(c echo.Context, message string, err ...error) error {
	var formattedError any
	if len(err) > 0 && err[0] != nil {
		formattedError = formatError(err[0])
	}

	return c.JSON(http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Message: message,
		Error:   formattedError,
	})
}

func NotFound(c echo.Context, message string, err ...error) error {
	var formattedError any
	if len(err) > 0 && err[0] != nil {
		formattedError = formatError(err[0])
	}

	return c.JSON(http.StatusNotFound, Response{
		Status:  http.StatusNotFound,
		Message: message,
		Error:   formattedError,
	})
}

func InternalServerError(c echo.Context, message string, err ...error) error {
	var formattedError any
	if len(err) > 0 && err[0] != nil {
		formattedError = formatError(err[0])
	}

	return c.JSON(http.StatusInternalServerError, Response{
		Status:  http.StatusInternalServerError,
		Message: message,
		Error:   formattedError,
	})
}

func Custom(c echo.Context, status int, message string, data any, err ...error) error {
	var formattedError any
	if len(err) > 0 && err[0] != nil {
		formattedError = formatError(err[0])
	}

	return c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   formattedError,
	})
}

func Forbidden(c echo.Context, message string, err ...error) error {
	var formattedError any
	if len(err) > 0 && err[0] != nil {
		formattedError = formatError(err[0])
	}

	return c.JSON(http.StatusForbidden, Response{
		Status:  http.StatusForbidden,
		Message: message,
		Error:   formattedError,
	})
}

func formatError(err error) any {
	if err == nil {
		return nil
	}
	return err.Error()
}
