package auth_cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
)

type AuthCmd struct {
	*core.Ctx
}

type InternalLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthCmd) InternalLogin(
	ctx context.Context,
	req InternalLoginRequest,
) (auth.AuthUser, auth.JWTToken, error) {
	authUser, cipherPassword, err := auth.AuthReader{}.FindUserByEmail(ctx, c.GormDB, req.Email)
	if err != nil {
		return auth.AuthUser{}, "", fmt.Errorf("InternalLogin: FindUserByEmail: %w", err)
	}

	if !auth.CheckCipherPassword(req.Password, cipherPassword) {
		return auth.AuthUser{}, "", errors.New("invalid password")
	}

	token, err := auth.GenerateJWTToken(c.JWTKey, authUser, auth.CreateTokenExpiry())
	if err != nil {
		return auth.AuthUser{}, "", fmt.Errorf("InternalLogin: GenerateJWTToken: %w", err)
	}

	return authUser, token, nil
}
