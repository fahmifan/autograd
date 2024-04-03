package user_management_cmd

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/user_management"
	"github.com/fahmifan/autograd/pkg/dbconn"
	"github.com/fahmifan/autograd/pkg/jobqueue/outbox"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	token, err := auth.GenerateRandomPlainPassword()
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateManagedUser: generate token")
		return nil, core.ErrInternalServer
	}

	now := time.Now()
	newUser, err := user_management.CreateManagedUser(user_management.CreateUserRequest{
		NewID:      uuid.New(),
		Now:        now,
		Name:       req.Msg.Name,
		Email:      req.Msg.Email,
		Role:       auth.Role(req.Msg.Role),
		NewTokenID: uuid.New(),
		Token:      token,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	password, err := auth.GenerateRandomPlainPassword()
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateManagedUser: GenerateRandomPassword")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// the password is a generated, user will have to change it, should be ok?
	cipherPassword, err := auth.WeakEncryptPassword(password)
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: CreateManagedUser: EncryptPassword")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	err = core.Transaction(ctx, cmd.Ctx, func(tx *gorm.DB) error {
		dbtx, err := dbconn.DBTxFromGorm(tx)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}

		err = user_management.ManagedUserWriter{}.SaveUserWithPasswordV2(ctx, dbtx, &newUser, cipherPassword)
		if err != nil {
			logs.ErrCtx(ctx, err, "UserManagementCmd: CreateManagedUser: SaveUserWithPassword")
			return connect.NewError(connect.CodeInternal, err)
		}

		_, err = cmd.OutboxEnqueuer.Enqueue(ctx, tx, outbox.EnqueueRequest{
			JobType: JobSendEmail,
			Payload: SendRegistrationEmailPayload{
				UserID: newUser.ID,
			},
		})
		if err != nil {
			logs.ErrCtx(ctx, err, "UserManagementCmd: CreateManagedUser: Enqueue")
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &connect.Response[autogradv1.CreatedResponse]{
		Msg: &autogradv1.CreatedResponse{
			Id:      newUser.ID.String(),
			Message: "user created",
		},
	}, nil
}

func (cmd *UserManagementCmd) ActivateManagedUser(
	ctx context.Context,
	req *connect.Request[autogradv1.ActivateManagedUserRequest],
) (*connect.Response[autogradv1.Empty], error) {
	userID, err := uuid.Parse(req.Msg.GetUserId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	now := time.Now()

	err = core.Transaction(ctx, cmd.Ctx, func(tx *gorm.DB) error {
		user, err := user_management.ManagedUserReader{}.FindUserByID(ctx, cmd.GormDB, userID)
		if err != nil {
			logs.ErrCtx(ctx, err, "UserManagementCmd: ActivateManagedUser: FindUserByID")
			return core.ErrInternalServer
		}

		user, err = user.Activate(now, req.Msg.GetActivationToken())
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = user_management.CheckPassword(req.Msg.GetPassword(), req.Msg.GetPassword())
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		cipherPassword, err := auth.EncryptPassword(req.Msg.GetPassword())
		if err != nil {
			logs.ErrCtx(ctx, err, "UserManagementCmd: ActivateManagedUser: EncryptPassword")
			return core.ErrInternalServer
		}

		err = user_management.ManagedUserWriter{}.SaveUserWithPassword(ctx, tx, false, user, cipherPassword)
		if err != nil {
			logs.ErrCtx(ctx, err, "UserManagementCmd: ActivateManagedUser: SaveUser")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &connect.Response[autogradv1.Empty]{}, nil
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
	token, err := auth.GenerateRandomPlainPassword()
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: InternalCreateAdminUser: generate token")
		return uuid.UUID{}, core.ErrInternalServer
	}

	newAdmin, err := user_management.CreateAdminUser(user_management.CreateAdminUserRequest{
		NewID:      uuid.New(),
		NewTokenID: uuid.New(),
		Token:      token,
		Now:        now,
		Email:      req.Email,
		Name:       req.Name,
	})
	if err != nil {
		return uuid.Nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	cipher, err := auth.EncryptPassword(req.Password)
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: InternalCreateAdminUser: EncryptPassword")
		return uuid.Nil, core.ErrInternalServer
	}

	err = user_management.ManagedUserWriter{}.SaveUserWithPassword(ctx, cmd.GormDB, true, newAdmin, cipher)
	if err != nil {
		logs.ErrCtx(ctx, err, "UserManagementCmd: InternalCreateAdminUser: SaveUserWithPassword")
		return uuid.Nil, core.ErrInternalServer
	}

	return newAdmin.ID, nil
}
