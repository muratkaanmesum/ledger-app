package customError

import (
	"errors"
	"fmt"
)

type HTTPStatusCode int

type Error struct {
	Code    HTTPStatusCode `json:"code"`
	Message string         `json:"message"`
	Details error          `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
}

func New(code HTTPStatusCode, message string, details ...error) *Error {
	var detail error
	if len(details) > 0 {
		detail = details[0]
	}

	return &Error{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

func BadRequest(message string, details ...error) *Error {
	return New(400, message, details...)
}

func Unauthorized(message string, details ...error) *Error {
	return New(401, message, details...)
}

func Forbidden(message string, details ...error) *Error {
	return New(403, message, details...)
}

func NotFound(message string, details ...error) *Error {
	return New(404, message, details...)
}

func InternalServerError(message string, details ...error) *Error {
	return New(500, message, details...)
}

func ServiceUnavailable(message string, details ...error) *Error {
	return New(503, message, details...)
}

func Parse(err error) *Error {
	if err == nil {
		return nil
	}

	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr
	}

	return InternalServerError("An unexpected error occurred", err)
}
