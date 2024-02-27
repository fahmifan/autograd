package httpsvc

import (
	"net/http"

	"connectrpc.com/connect"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/labstack/echo/v4"
)

func (s *Server) handleActivateManagedUser(c echo.Context) error {
	token := c.QueryParam("activationToken")
	userID := c.QueryParam("userID")

	_, err := s.service.ActivateManagedUser(c.Request().Context(), &connect.Request[autogradv1.ActivateManagedUserRequest]{
		Msg: &autogradv1.ActivateManagedUserRequest{
			UserId:          userID,
			ActivationToken: token,
		},
	})
	if err != nil {
		return responseConnectError(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}
