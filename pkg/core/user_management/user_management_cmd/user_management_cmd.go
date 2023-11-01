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

func (cmd *UserManagementCmd) CreateManagedUser(
	ctx context.Context,
	req *connect.Request[autogradv1.CreateManagedUserRequest],
) (*connect.Response[autogradv1.CreatedResponse], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.CreateAnyUser) {
		return nil, core.ErrPermissionDenied
	}

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

	password, err := auth.GenerateRandomPlainPassword()
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateUser: GenerateRandomPassword")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	cipherPassword, err := auth.EncryptPassword(password)
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateUser: EncryptPassword")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	err = user_management.ManagedUserWriter{}.SaveUserWithPassword(ctx, cmd.GormDB, newUser, cipherPassword)
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

type CreateAdminUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cmd *UserManagementCmd) InternalCreateAdminUser(
	ctx context.Context,
	req CreateAdminUserRequest,
) (id uuid.UUID, err error) {
	now := time.Now()
	newAdmin, err := user_management.CreateAdminUser(user_management.CreateAdminUserRequest{
		NewID: uuid.New(),
		Now:   now,
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		return uuid.Nil, err
	}

	cipher, err := auth.EncryptPassword(req.Password)
	if err != nil {
		return uuid.Nil, err
	}

	err = user_management.ManagedUserWriter{}.SaveUserWithPassword(ctx, cmd.GormDB, newAdmin, cipher)
	if err != nil {
		return uuid.Nil, err
	}

	return newAdmin.ID, nil
}
