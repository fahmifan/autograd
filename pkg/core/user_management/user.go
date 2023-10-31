package user_management

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type ManagedUser struct {
	ID     uuid.UUID
	Name   string
	Email  string
	Role   auth.Role
	Active bool
	core.TimestampMetadata
}

type CreateUserRequest struct {
	NewID uuid.UUID
	Now   time.Time
	Name  string
	Email string
	Role  auth.Role
}

func CreateUser(req CreateUserRequest) (ManagedUser, error) {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ManagedUser{}, ErrInvalidEmail
	}

	if !auth.ValidRole(req.Role) {
		return ManagedUser{}, errors.New("invalid role")
	}

	if len(strings.TrimSpace(req.Name)) < 3 {
		return ManagedUser{}, errors.New("name must be at least 3 characters long")
	}

	return ManagedUser{
		ID:                req.NewID,
		Name:              req.Name,
		Email:             req.Email,
		Role:              req.Role,
		TimestampMetadata: core.NewEntityMeta(req.Now),
		Active:            true,
	}, nil
}

type CreateAdminUserRequest struct {
	NewID uuid.UUID
	Now   time.Time
	Name  string
	Email string
}

func CreateAdminUser(req CreateAdminUserRequest) (ManagedUser, error) {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ManagedUser{}, ErrInvalidEmail
	}

	if len(strings.TrimSpace(req.Name)) < 3 {
		return ManagedUser{}, errors.New("name must be at least 3 characters long")
	}

	return ManagedUser{
		ID:                req.NewID,
		Name:              req.Name,
		Email:             req.Email,
		Role:              auth.RoleAdmin,
		TimestampMetadata: core.NewEntityMeta(req.Now),
		Active:            true,
	}, nil
}
