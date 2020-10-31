package httpsvc

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
)

type createUserReq struct {
	Email    string
	Password string
}

type userRes struct {
	ID        string
	Email     string
	Role      string
	CreatedAt string
	UpdatedAt string
}

func userResFromModel(m *model.User) *userRes {
	return &userRes{
		ID:        utils.Int64ToString(m.ID),
		Email:     m.Email,
		Role:      m.Role.ToString(),
		CreatedAt: m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339Nano),
	}
}

func responseError(c echo.Context, err error) error {
	switch err {
	case nil:
		return c.JSON(http.StatusOK, nil)
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) handleCreateUser(c echo.Context) error {
	user := &model.User{}

	c.Bind(user)
	err := s.userUsecase.Create(c.Request().Context(), user)
	if err != nil {
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, userResFromModel(user))
}
