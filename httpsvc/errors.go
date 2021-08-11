package httpsvc

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// errors ..
var (
	ErrNotFound = errors.New("not found")
)

func responseError(c echo.Context, err error) error {
	switch err {
	case nil:
		return c.JSON(http.StatusOK, nil)
	case ErrUnauthorized:
		return c.JSON(http.StatusUnauthorized, Error{Error: "unauthorized"})
	case ErrNotFound:
		return c.JSON(http.StatusNotFound, Error{Error: "not found"})
	default:
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, Error{Error: "system error"})
	}
}
