package model

import (
	"time"
)

// UserRole :nodoc:
type UserRole int

// user roles
const (
	UserRoleAdmin   = UserRole(1)
	UserRoleStudent = UserRole(2)
)

// ToString :nodoc:
func (u UserRole) ToString() string {
	switch u {
	case UserRoleAdmin:
		return "ADMIN"
	case UserRoleStudent:
		return "STUDENT"
	default:
		return ""
	}
}

// User :nodoc:
type User struct {
	ID        int64
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
