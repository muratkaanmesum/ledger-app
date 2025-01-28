package middlewares

import (
	"github.com/labstack/echo/v4"
	"ptm/pkg/counter"
)

func PerformanceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			counter.AddTotalRequests()
			err := next(c)
			if err != nil {
				counter.AddFail()
				return err
			}

			if c.Response().Status < 400 {
				counter.AddSuccess()
			}

			return nil
		}
	}
}
