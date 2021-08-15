package httpsvc

import (
	"net/http"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const accessTokenExpiry = 30 * time.Minute

// Login godoc
// @Summary login a User
// @Description currently it only support one session per user
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
	ctx := c.Request().Context()

	user, err := s.userUsecase.FindByEmailAndPassword(ctx, req.Email, req.Password)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	err = s.sessionRepo.DeleteAllByUserID(ctx, user.ID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	// generate session & refresh token
	now := time.Now()
	const sevenDay = time.Hour * 24 * 7
	sess := &model.Session{
		UserID:    user.ID,
		ExpiredAt: now.Add(sevenDay),
	}
	err = s.sessionRepo.Create(ctx, sess)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	expiredAt := now.Add(accessTokenExpiry)
	accessToken, err := generateAccessToken(user, expiredAt)
	if err != nil {
		logrus.WithField("email", req.Email).Error(err)
		return responseError(c, err)
	}

	refreshToken, err := generateRefreshToken(sess)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiredIn:    int32(expiredAt.Unix()),
	})
}

func (s *Server) handleRefreshToken(c echo.Context) error {
	req := struct {
		RT string `json:"refreshToken"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		logrus.Error(err)
		return responseError(c, ErrInvalidArguments)
	}

	sessID, err := parseRefreshToken(req.RT)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	ctx := c.Request().Context()

	oldSess, err := s.sessionRepo.FindByID(ctx, sessID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}
	if oldSess == nil {
		return responseError(c, ErrNotFound)
	}

	now := time.Now()
	isExpired := now.After(oldSess.ExpiredAt)
	if isExpired {
		return c.JSON(http.StatusBadRequest, Error{Error: "token already expired"})
	}

	user, err := s.userUsecase.FindByID(ctx, oldSess.UserID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	expiredAt := now.Add(accessTokenExpiry)
	accessToken, err := generateAccessToken(user, expiredAt)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RT,
		TokenType:    "Bearer",
		ExpiredIn:    int32(expiredAt.Unix()),
	})
}
