package httpsvc

import (
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/labstack/echo/v4"
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
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	case ErrNotFound:
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	default:
		logs.ErrCtx(c.Request().Context(), err, "responseError")
		return c.JSON(http.StatusInternalServerError, "system error")
	}
}

func responseConnectError(ec echo.Context, err error) error {
	connectErr, ok := err.(*connect.Error)
	if !ok {
		return responseError(ec, err)
	}

	code := connectErr.Code()

	switch code {
	case connect.CodeInternal:
		return ec.JSON(http.StatusInternalServerError, connectErr)
	default:
		return ec.JSON(http.StatusBadRequest, echo.Map{
			"error": connectErr.Error(),
		})
	}
}
