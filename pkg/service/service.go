package service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/assignments/assignments_cmd"
	"github.com/fahmifan/autograd/pkg/core/assignments/assignments_query"
	"github.com/fahmifan/autograd/pkg/core/auth/auth_cmd"
	"github.com/fahmifan/autograd/pkg/core/mediastore/mediastore_cmd"
	"github.com/fahmifan/autograd/pkg/core/student_assignment/student_assignment_cmd"
	"github.com/fahmifan/autograd/pkg/core/student_assignment/student_assignment_query"
	"github.com/fahmifan/autograd/pkg/core/user_management/user_management_cmd"
	"github.com/fahmifan/autograd/pkg/core/user_management/user_management_query"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/jobqueue/outbox"
	"github.com/fahmifan/autograd/pkg/mailer"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"gorm.io/gorm"
)

type Service struct {
	coreCtx *core.Ctx

	*auth_cmd.AuthCmd
	*user_management_query.UserManagementQuery
	*user_management_cmd.UserManagementCmd
	*assignments_cmd.AssignmentCmd
	*assignments_query.AssignmentsQuery
	*mediastore_cmd.MediaStoreCmd
	*student_assignment_query.StudentAssignmentQuery
	*student_assignment_cmd.StudentAssignmentCmd

	outboxService *outbox.OutboxService
}

var _ autogradv1connect.AutogradServiceHandler = &Service{}

func NewService(
	gormDB *gorm.DB,
	jwtKey string,
	mediaCfg core.MediaConfig,
	senderEmail string,
	mailer mailer.Mailer,
) *Service {
	outboxService := outbox.NewOutboxService(gormDB)

	coreCtx := &core.Ctx{
		GormDB:         gormDB,
		JWTKey:         jwtKey,
		MediaConfig:    mediaCfg,
		SenderEmail:    senderEmail,
		AppLink:        "http://localhost:5173",
		LogoURL:        "http://localhost:5173/logo.png",
		Mailer:         mailer,
		OutboxEnqueuer: outboxService,
	}

	return &Service{
		coreCtx:                coreCtx,
		AuthCmd:                &auth_cmd.AuthCmd{Ctx: coreCtx},
		UserManagementCmd:      &user_management_cmd.UserManagementCmd{Ctx: coreCtx},
		UserManagementQuery:    &user_management_query.UserManagementQuery{Ctx: coreCtx},
		AssignmentCmd:          &assignments_cmd.AssignmentCmd{Ctx: coreCtx},
		AssignmentsQuery:       &assignments_query.AssignmentsQuery{Ctx: coreCtx},
		MediaStoreCmd:          &mediastore_cmd.MediaStoreCmd{Ctx: coreCtx},
		StudentAssignmentQuery: &student_assignment_query.StudentAssignmentQuery{Ctx: coreCtx},
		StudentAssignmentCmd:   &student_assignment_cmd.StudentAssignmentCmd{Ctx: coreCtx},
		outboxService:          outboxService,
	}
}

func (service *Service) Ping(ctx context.Context, req *connect.Request[autogradv1.Empty]) (*connect.Response[autogradv1.PingResponse], error) {
	return &connect.Response[autogradv1.PingResponse]{
		Msg: &autogradv1.PingResponse{
			Message: "pong",
		},
	}, nil
}

func (service *Service) RunOutboxService(ctx context.Context) error {
	return service.outboxService.Run(ctx)
}

func (service *Service) RegisterJobHandlers() {
	handlers := []jobqueue.JobHandler{
		&user_management_cmd.SendRegistrationEmailHandler{Ctx: service.coreCtx},
	}

	outbox.RegisterHandlers(service.coreCtx.GormDB, handlers)
}
