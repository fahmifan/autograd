package service

import (
	"context"
	"database/sql"

	"connectrpc.com/connect"
	"github.com/coocood/freecache"
	"github.com/fahmifan/autograd/pkg/config"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/admin_courses/admin_courses_cmd"
	"github.com/fahmifan/autograd/pkg/core/admin_courses/admin_courses_query"
	"github.com/fahmifan/autograd/pkg/core/assignments/assignments_cmd"
	"github.com/fahmifan/autograd/pkg/core/assignments/assignments_query"
	"github.com/fahmifan/autograd/pkg/core/auth/auth_cmd"
	"github.com/fahmifan/autograd/pkg/core/grading/grading_cmd"
	"github.com/fahmifan/autograd/pkg/core/mediastore/mediastore_cmd"
	"github.com/fahmifan/autograd/pkg/core/student_assignment/student_assignment_cmd"
	"github.com/fahmifan/autograd/pkg/core/student_assignment/student_assignment_query"
	"github.com/fahmifan/autograd/pkg/core/student_courses/student_courses_cmdquery"
	"github.com/fahmifan/autograd/pkg/core/user_management/user_management_cmd"
	"github.com/fahmifan/autograd/pkg/core/user_management/user_management_query"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/jobqueue/outbox"
	"github.com/fahmifan/autograd/pkg/mailer"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"golang.org/x/sync/singleflight"
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
	*grading_cmd.GradingCmd
	*admin_courses_query.AdminCoursesQuery
	*admin_courses_cmd.AdminCoursesCmd
	*student_courses_cmdquery.StudentCoursesQuery

	outboxService *outbox.OutboxService
}

var _ autogradv1connect.AutogradServiceHandler = &Service{}

func NewService(
	gormDB *gorm.DB,
	sqlDB *sql.DB,
	jwtKey string,
	debug bool,
	mediaCfg core.MediaConfig,
	senderEmail string,
	mailer mailer.Mailer,
	cache *freecache.Cache,
) *Service {
	outboxService := outbox.NewOutboxService(gormDB, sqlDB, config.Debug())

	coreCtx := &core.Ctx{
		GormDB:         gormDB,
		JWTKey:         jwtKey,
		MediaConfig:    mediaCfg,
		SenderEmail:    senderEmail,
		AppLink:        config.WebBaseURL(),
		LogoURL:        config.BaseURL() + "/logo.png",
		Mailer:         mailer,
		OutboxEnqueuer: outboxService,
		SqlDB:          sqlDB,
		Debug:          debug,
		Cache:          cache,
		Flight:         &singleflight.Group{},
	}

	return &Service{
		coreCtx:                coreCtx,
		outboxService:          outboxService,
		AuthCmd:                &auth_cmd.AuthCmd{Ctx: coreCtx},
		UserManagementCmd:      &user_management_cmd.UserManagementCmd{Ctx: coreCtx},
		UserManagementQuery:    &user_management_query.UserManagementQuery{Ctx: coreCtx},
		AssignmentCmd:          &assignments_cmd.AssignmentCmd{Ctx: coreCtx},
		AssignmentsQuery:       &assignments_query.AssignmentsQuery{Ctx: coreCtx},
		MediaStoreCmd:          &mediastore_cmd.MediaStoreCmd{Ctx: coreCtx},
		StudentAssignmentQuery: &student_assignment_query.StudentAssignmentQuery{Ctx: coreCtx},
		StudentAssignmentCmd:   &student_assignment_cmd.StudentAssignmentCmd{Ctx: coreCtx},
		GradingCmd:             &grading_cmd.GradingCmd{Ctx: coreCtx},
		AdminCoursesQuery:      &admin_courses_query.AdminCoursesQuery{Ctx: coreCtx},
		AdminCoursesCmd:        &admin_courses_cmd.AdminCoursesCmd{Ctx: coreCtx},
		StudentCoursesQuery:    &student_courses_cmdquery.StudentCoursesQuery{Ctx: coreCtx},
	}
}

func (service *Service) Ping(ctx context.Context, req *connect.Request[autogradv1.Empty]) (*connect.Response[autogradv1.PingResponse], error) {
	return &connect.Response[autogradv1.PingResponse]{
		Msg: &autogradv1.PingResponse{
			Message: "pong",
		},
	}, nil
}

func (service *Service) RunOutboxService() error {
	return service.outboxService.Run()
}

func (service *Service) StopOutboxService() {
	service.outboxService.Stop()
}

func (service *Service) RegisterJobHandlers() {
	handlers := []jobqueue.JobHandler{
		&user_management_cmd.SendRegistrationEmailHandler{Ctx: service.coreCtx},
		&student_assignment_cmd.GradeStudentSubmissionHandler{Ctx: service.coreCtx},
	}

	outbox.RegisterHandlers(service.coreCtx.GormDB, service.coreCtx.SqlDB, service.coreCtx.Debug, handlers)
}
