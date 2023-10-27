package httpsvc

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ErrUnauthorized error
var ErrUnauthorized = errors.New("unauthorized")

// ErrMissingAuthorization error
var ErrMissingAuthorization = errors.New("missing Authorization header")

func (server *Server) addUserToCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := parseTokenFromHeader(&c.Request().Header)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid token"})
		}

		user, ok := parseToken(token)
		if !ok || user == nil {
			return next(c)
		}

		setUserToCtx(c, user)
		userID, err := uuid.Parse(user.ID)

		ctx := c.Request().Context()
		authUser := auth.AuthUser{
			UserID: userID,
			Role:   auth.Role(string(user.Role)),
		}
		reqCtx := auth.CtxWithUser(ctx, authUser)
		c.SetRequest(c.Request().WithContext(reqCtx))

		return next(c)
	}
}

func (s *Server) authz(perms ...auth.Permission) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := getUserFromCtx(c)
			if user == nil {
				log.Error(ErrUnauthorized)
				return responseError(c, ErrUnauthorized)
			}

			for _, p := range perms {
				if user.Role.Granted(p) {
					return next(c)
				}
			}

			return responseError(c, ErrUnauthorized)
		}
	}
}

const authzHeader = "Authorization"

func parseTokenFromHeader(header *http.Header) (string, error) {
	var token string

	authHeaders := strings.Split(header.Get(authzHeader), " ")
	if len(authHeaders) != 2 {
		return "", ErrTokenInvalid
	}

	if authHeaders[0] != "Bearer" {
		err := ErrMissingAuthorization
		log.WithField(authzHeader, header.Get(authzHeader)).Error(err)
		return token, err
	}

	token = strings.Trim(authHeaders[1], " ")
	if token == "" {
		err := ErrMissingAuthorization
		log.WithField(authzHeader, header.Get(authzHeader)).Error(err)
		return token, err
	}

	return token, nil
}
