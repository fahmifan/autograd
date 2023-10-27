package user_management_cmd

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/user_management"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/google/uuid"
)

type UserManagementCmd struct {
	*core.Ctx
}

func (cmd *UserManagementCmd) CreateUser(
	ctx context.Context,
	req *connect.Request[autogradv1.CreateUserRequest],
) (*connect.Response[autogradv1.CreatedResponse], error) {
	now := time.Now()
	newUser, err := user_management.CreateUser(user_management.CreateUserRequest{
		NewID: uuid.New(),
		Now:   now,
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
		Role:  auth.Role(req.Msg.Role),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	password, err := user_management.GenerateRandomPassword()
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateUser: GenerateRandomPassword")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	err = user_management.UserWriter{}.SaveUserWithPassword(ctx, cmd.GormDB, newUser, password)
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateUser: SaveUserWithPassword")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[autogradv1.CreatedResponse]{
		Msg: &autogradv1.CreatedResponse{
			Id:      newUser.ID.String(),
			Message: "user created",
		},
	}, nil
}
