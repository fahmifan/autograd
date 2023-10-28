package httpsvc

import (
	"net/http"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const userInfoCtx = "userInfoCtx"

type userRequest struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     auth.Role `json:"role"`
}

func (u *userRequest) toModel() *model.User {
	return &model.User{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Role:     u.Role,
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
		ID:        m.ID,
		Email:     m.Email,
		Role:      m.Role.ToString(),
		CreatedAt: m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339Nano),
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

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		logrus.WithField("uuid", user.ID).Error(err)
		return responseError(c, err)
	}

	token, err := auth.GenerateJWTToken(
		s.jwtKey,
		auth.AuthUser{
			UserID: userID,
			Email:  user.Email,
			Role:   user.Role,
			Name:   user.Name,
		}, auth.CreateTokenExpiry())
	if err != nil {
		logrus.WithField("email", userReq.Email).Error(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
