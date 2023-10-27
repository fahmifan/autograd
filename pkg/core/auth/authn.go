package auth

import (
	"context"

	"github.com/google/uuid"
)

type CtxKey string

const userInfoCtx CtxKey = "core.auth.userInfoCtx"

type AuthUser struct {
	UserID uuid.UUID
	Role   Role
}

func GetUserFromCtx(ctx context.Context) (AuthUser, bool) {
	val, ok := ctx.Value(userInfoCtx).(AuthUser)
	return val, ok
}

func CtxWithUser(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, userInfoCtx, user)
}
