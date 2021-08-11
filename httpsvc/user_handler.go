package httpsvc

import (
	"net/http"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const userInfoCtx = "userInfoCtx"

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role" enums:"ADMIN,STUDENT"`
}

func (u *UserRequest) toModel() *model.User {
	return &model.User{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Role:     model.ParseRole(u.Role),
	}
}

type UserRes struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func userResFromModel(m *model.User) *UserRes {
	return &UserRes{
		ID:        m.ID,
		Email:     m.Email,
		Role:      m.Role.ToString(),
		CreatedAt: m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339Nano),
	}
}

type Error struct {
	Error string `json:"error"`
}

// CreateUser godoc
// @Summary create (register) a User
// @Description create a User
// @ID CreateUser
// @Accept  json
// @Produce  json
// @Param user body UserRequest true "name"
// @Success 200 {object} UserRes
// @Failure default {object} Error
// @Router /api/v1/users [post]
func (s *Server) handleCreateUser(c echo.Context) error {
	userReq := &UserRequest{}
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

type LoginRequest struct {
	Email    string `json:"email" example:"your@email.com"`
	Password string `json:"password"`
} //@name LoginRequest

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiredIn    int32  `json:"expiredIn"`
	RefreshToken string `json:"refreshToken"`
}

// Login godoc
// @Summary login a User
// @ID Login
// @Accept json
// @Produce json
// @Param user body LoginRequest true "login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} Error
// @Router /api/v1/users/login [post]
func (s *Server) handleLogin(c echo.Context) error {
	req := &LoginRequest{}
	err := c.Bind(req)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	user, err := s.userUsecase.FindByEmailAndPassword(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	expiredIn := createTokenExpiry()
	token, err := generateToken(*user, expiredIn)
	if err != nil {
		logrus.WithField("email", req.Email).Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiredIn:   int32(expiredIn),
	})
}
