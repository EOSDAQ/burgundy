package middleware

import (
	"burgundy/util"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// TransID Transaction ID(yyyymmddhhmi + 5 numbers)
func TransID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(
		middleware.RequestIDConfig{
			Generator: util.NewID,
		})
}
