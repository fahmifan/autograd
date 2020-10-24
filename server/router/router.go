package router

import (
	"autograd/server/api"
	"autograd/server/api/middleware"

	"github.com/labstack/echo"
)

func New() *echo.Echo {

	e := echo.New()

	middleware.SetMiddleware(e)
	api.MainGroup(e)

	return e
}
