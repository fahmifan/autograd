package auth

import (
	"context"

	"github.com/google/uuid"
	passwordgen "github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

type CtxKey string

const userInfoCtx CtxKey = "core.auth.userInfoCtx"

type AuthUser struct {
	UserID uuid.UUID
	Email  string
	Name   string
	Role   Role
}

func GetUserFromCtx(ctx context.Context) (AuthUser, bool) {
	val, ok := ctx.Value(userInfoCtx).(AuthUser)
	return val, ok
}

func CtxWithUser(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, userInfoCtx, user)
}

type CipherPassword string

func EncryptPassword(password string) (CipherPassword, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return CipherPassword(bytes), err
}

func CheckCipherPassword(password string, hash CipherPassword) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomPlainPassword() (string, error) {
	pass, err := passwordgen.Generate(12, 3, 2, false, false)
	return pass, err
}
