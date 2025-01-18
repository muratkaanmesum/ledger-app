package customError

import "fmt"

type HTTPStatusCode int

const (
	BadRequest          HTTPStatusCode = 400
	Unauthorized        HTTPStatusCode = 401
	Forbidden           HTTPStatusCode = 403
	NotFound            HTTPStatusCode = 404
	InternalServerError HTTPStatusCode = 500
	ServiceUnavailable  HTTPStatusCode = 503
)

type Error struct {
	Code HTTPStatusCode
	err  error
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%d] %s: %s - %v", e.Code, e.err)
	}
	return fmt.Sprintf("[%d] %s: %s", e.Code)
}

func (e *Error) Unwrap() error {
	return e.err
}

func New(code HTTPStatusCode, err ...error) *Error {
	var wrappedErr error
	if len(err) > 0 && err[0] != nil {
		wrappedErr = err[0]
	}
	return &Error{
		Code: code,
		err:  wrappedErr,
	}
}

func Wrap(code HTTPStatusCode, err ...error) *Error {
	var wrappedErr error
	if len(err) > 0 && err[0] != nil {
		wrappedErr = err[0]
	}
	return &Error{
		Code: code,
		err:  wrappedErr,
	}
}
