package api

import (
	"autograd/server/handler"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/test", handler.Test)
	e.POST("/upload", handler.Upload)
}
