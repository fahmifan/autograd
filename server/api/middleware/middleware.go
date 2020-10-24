package middleware

import (
	"github.com/labstack/echo"
)

func SetMiddleware(e *echo.Echo) {
	e.Use(ServerHeader)
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Custom-Header", "this is a custom header")
		return next(c)
	}
}
