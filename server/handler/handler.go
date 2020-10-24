package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func Test(context echo.Context) error {
	return context.String(http.StatusOK, fmt.Sprintf("this is a test"))
}
