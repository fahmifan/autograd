package httpsvc

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

const userInfoCtx = "userInfoCtx"

type userRequest struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func (u *userRequest) toModel() *model.User {
	return &model.User{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Role:     model.ParseRole(u.Role),
	}
}

type userRes struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
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
	userReq := &userRequest{}
	err := c.Bind(userReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	user := userReq.toModel()
	err = s.userUsecase.Create(c.Request().Context(), user)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, userResFromModel(user))
}

func (s *Server) handleLogin(c echo.Context) error {
	userReq := &userRequest{}
	err := c.Bind(userReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	user, err := s.userUsecase.FindByEmailAndPassword(c.Request().Context(), userReq.Email, userReq.Password)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	token, err := generateToken(*user, createTokenExpiry())
	if err != nil {
		logrus.WithField("email", userReq.Email).Error(err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
