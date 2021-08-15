package httpsvc

import (
	"errors"
	"net/http"

	"github.com/fahmifan/autograd/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// errors ..
var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidArguments = errors.New("invalid arguments")
)

func responseError(c echo.Context, err error) error {
	switch err {
	case nil:
		return c.JSON(http.StatusOK, nil)
	case ErrUnauthorized:
		return c.JSON(http.StatusUnauthorized, Error{Error: "unauthorized"})
	case ErrNotFound, usecase.ErrNotFound:
		return c.JSON(http.StatusNotFound, Error{Error: "not found"})
	case ErrTokenInvalid:
		return c.JSON(http.StatusUnauthorized, Error{Error: "unauthorized"})
	case ErrInvalidArguments, usecase.ErrInvalidArguments:
		return c.JSON(http.StatusBadRequest, Error{Error: "invalid arguments"})
	default:
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, Error{Error: "system error"})
	}
}
