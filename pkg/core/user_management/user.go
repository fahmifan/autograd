package user_management

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/google/uuid"
	passwordgen "github.com/sethvargo/go-password/password"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
	Role  auth.Role
	core.EntityMeta
}

type CreateUserRequest struct {
	NewID uuid.UUID
	Now   time.Time
	Name  string
	Email string
	Role  auth.Role
}

func CreateUser(req CreateUserRequest) (User, error) {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return User{}, ErrInvalidEmail
	}

	if !auth.ValidRole(req.Role) {
		return User{}, errors.New("invalid role")
	}

	if len(strings.TrimSpace(req.Name)) < 3 {
		return User{}, errors.New("name must be at least 3 characters long")
	}

	return User{
		ID:         req.NewID,
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		EntityMeta: core.NewEntityMeta(req.Now),
	}, nil
}

func GenerateRandomPassword() (string, error) {
	return passwordgen.Generate(12, 8, 4, false, false)
}
