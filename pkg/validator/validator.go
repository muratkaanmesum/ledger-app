package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// New creates a new CustomValidator
func New() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate method for Echo to use
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func RegisterValidator(e *echo.Echo) {
	e.Validator = New()
}
