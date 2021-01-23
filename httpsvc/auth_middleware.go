package httpsvc

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fahmifan/autograd/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ErrUnauthorized error
var ErrUnauthorized = errors.New("unauthorized")

// ErrMissingAuthorization error
var ErrMissingAuthorization = errors.New("missing Authorization header")

// authorizeByRoleMiddleware authorized request by given authorized roles
func (s *Server) authorizeByRoleMiddleware(authorizedRole []model.Role) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := getUserFromCtx(c)
			if user == nil {
				log.Warn("user nil")
				return responseError(c, ErrUnauthorized)
			}

			if ok := authorizeByRole(user.Role, authorizedRole); !ok {
				log.WithField("role", user.Role).Warn("unauthorized role")
				return responseError(c, ErrUnauthorized)
			}

			return next(c)
		}
	}
}

func authorizeByRole(userRole model.Role, authorizedRole []model.Role) bool {
	for _, role := range authorizedRole {
		if userRole == role {
			return true
		}
	}

	return false
}

// AuthMiddleware ..
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := parseTokenFromHeader(&c.Request().Header)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		user, ok := auth(token)
		if !ok {
			log.Error(ErrUnauthorized)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": ErrUnauthorized.Error()})
		}

		setUserToCtx(c, user)
		return next(c)
	}
}

func parseTokenFromHeader(header *http.Header) (string, error) {
	var token string

	authHeaders := strings.Split(header.Get("Authorization"), " ")
	if authHeaders[0] != "Bearer" {
		err := ErrMissingAuthorization
		log.WithField("Authorization", header.Get("Authorization")).Error(err)
		return token, err
	}

	token = strings.Trim(authHeaders[1], " ")
	if token == "" {
		err := ErrMissingAuthorization
		log.WithField("Authorization", header.Get("Authorization")).Error(err)
		return token, err
	}

	return token, nil
}
