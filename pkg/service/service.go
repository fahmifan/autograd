package service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/user_management/user_management_cmd"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"gorm.io/gorm"
)

type Service struct {
	*user_management_cmd.UserManagementCmd
}

var _ autogradv1connect.AutogradServiceHandler = &Service{}

func NewService(gormDB *gorm.DB) *Service {
	coreCtx := &core.Ctx{
		GormDB: gormDB,
	}
	return &Service{
		UserManagementCmd: &user_management_cmd.UserManagementCmd{Ctx: coreCtx},
	}
}

func (service *Service) Ping(ctx context.Context, req *connect.Request[autogradv1.Empty]) (*connect.Response[autogradv1.PingResponse], error) {
	return &connect.Response[autogradv1.PingResponse]{
		Msg: &autogradv1.PingResponse{
			Message: "pong",
		},
	}, nil
}