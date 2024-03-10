package auth_cmd

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
)

type AuthCmd struct {
	*core.Ctx
}

type InternalLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cmd *AuthCmd) InternalLogin(
	ctx context.Context,
	req InternalLoginRequest,
) (auth.AuthUser, auth.JWTToken, error) {
	authUser, cipherPassword, err := auth.AuthReader{}.FindUserByEmail(ctx, cmd.GormDB, req.Email)
	if err != nil {
		return auth.AuthUser{}, "", fmt.Errorf("InternalLogin: FindUserByEmail: %w", err)
	}

	if !auth.CheckCipherPassword(req.Password, cipherPassword) {
		return auth.AuthUser{}, "", errors.New("invalid password")
	}

	token, err := auth.GenerateJWTToken(cmd.JWTKey, authUser, auth.CreateTokenExpiry())
	if err != nil {
		return auth.AuthUser{}, "", fmt.Errorf("InternalLogin: GenerateJWTToken: %w", err)
	}

	return authUser, token, nil
}

func (cmd *AuthCmd) Login(ctx context.Context, req *connect.Request[autogradv1.LoginRequest]) (*connect.Response[autogradv1.LoginResponse], error) {
	authUser, cipherPassword, err := auth.AuthReader{}.FindUserByEmail(ctx, cmd.GormDB, req.Msg.GetEmail())
	if err != nil {
		logs.ErrCtx(ctx, err, "AuthCmd: Login: FindUserByEmail")
		return nil, core.ErrInternalServer
	}

	if !auth.CheckCipherPassword(req.Msg.GetPassword(), cipherPassword) {
		cerr := connect.NewError(connect.CodeInvalidArgument, err)
		return nil, cerr
	}

	token, err := auth.GenerateJWTToken(cmd.JWTKey, authUser, auth.CreateTokenExpiry())
	if err != nil {
		logs.ErrCtx(ctx, err, "AuthCmd: Login: GenerateJWTToken")
		return nil, core.ErrInternalServer
	}

	return &connect.Response[autogradv1.LoginResponse]{
		Msg: &autogradv1.LoginResponse{
			Token: string(token),
		},
	}, nil
}
