package httpsvc

import (
	"context"
	"net/http"
	"strings"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"github.com/fahmifan/autograd/pkg/service"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// Server ..
type Server struct {
	echo              *echo.Echo
	gormDB            *gorm.DB
	port              string
	staticMediaPath   string
	userUsecase       model.UserUsecase
	assignmentUsecase model.AssignmentUsecase
	submissionUsecase model.SubmissionUsecase
	mediaUsecase      model.MediaUsecase
	service           *service.Service
	jwtKey            string
}

// NewServer ..
func NewServer(port, staticMediaPath string, opts ...Option) *Server {
	s := &Server{
		echo:            echo.New(),
		port:            port,
		staticMediaPath: staticMediaPath,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Run server
func (s *Server) Run() {
	s.routes()
	logrus.Fatal(s.echo.Start(":" + s.port))
}

// Stop server gracefully
func (s *Server) Stop(ctx context.Context) {
	if err := s.echo.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}
}

func (s *Server) routes() {
	s.echo.Use(
		middleware.CORS(),
		s.addUserToCtx,
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: logs.EchoRequestLogger(true),
			LogLatency:    true,
			LogRemoteIP:   true,
			LogUserAgent:  true,
			LogError:      true,
			HandleError:   true,
		}),
		// FIXME: debug mode
	)

	s.echo.GET("/ping", s.handlePing)

	// TODO: add auth for private static
	s.echo.Static("/storage", "submission")
	s.echo.Static("/media", s.staticMediaPath)

	apiV1 := s.echo.Group("/api/v1")
	apiV1.POST("/users", s.handleCreateUser)
	apiV1.POST("/users/login", s.handleLogin)

	apiV1.POST("/assignments", s.handleCreateAssignment, s.authz(auth.CreateAssignment))
	apiV1.GET("/assignments", s.handleGetAssignments, s.authz(auth.ViewAnyAssignments))
	apiV1.GET("/assignments/:id", s.handleGetAssignment, s.authz(auth.ViewAssignment))
	apiV1.GET("/assignments/:id/submissions", s.handleGetAssignmentSubmissions, s.authz(auth.ViewAnySubmissions))
	apiV1.PUT("/assignments/:id", s.handleUpdateAssignment, s.authz(auth.UpdateAssignment))
	apiV1.DELETE("/assignments/:id", s.handleDeleteAssignment, s.authz(auth.DeleteAssignment))

	apiV1.POST("/submissions", s.handleCreateSubmission, s.authz(auth.CreateSubmission))
	apiV1.GET("/submissions/:id", s.handleGetSubmission, s.authz(auth.ViewAnySubmissions))
	apiV1.PUT("/submissions", s.handleUpdateSubmission, s.authz(auth.UpdateSubmission))
	apiV1.DELETE("/submissions/:id", s.handleDeleteSubmission, s.authz(auth.DeleteSubmission))

	apiV1.POST("/media/upload", s.handleUploadMedia, s.authz(auth.CreateMedia))

	grpHandlerName, grpcHandler := autogradv1connect.NewAutogradServiceHandler(
		s.service,
	)
	s.echo.Group("/grpc").Any(
		grpHandlerName+"*",
		echo.WrapHandler(grpcHandler),
		trimPathGroup("/grpc"),
	)
}

func trimPathGroup(groupPrefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().URL.Path = strings.TrimPrefix(c.Request().URL.Path, groupPrefix)
			return next(c)
		}
	}
}

func (s *Server) handlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
