package httpsvc

import (
	"context"
	"net/http"
	"strings"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"github.com/fahmifan/autograd/pkg/service"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
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
	s.echo.GET("/ping", s.handlePing)

	// TODO: add auth for private static
	s.echo.Static("/storage", "submission")
	s.echo.Static("/media", s.staticMediaPath)

	apiV1 := s.echo.Group("/api/v1")
	apiV1.POST("/users", s.handleCreateUser)
	apiV1.POST("/users/login", s.handleLogin)

	apiV1.POST("/assignments", s.handleCreateAssignment, s.authorizedOne(model.CreateAssignment))
	apiV1.GET("/assignments", s.handleGetAssignments, s.authorizedOne(model.ViewAnyAssignments))
	apiV1.GET("/assignments/:id", s.handleGetAssignment, s.authorizedOne(model.ViewAssignment))
	apiV1.GET("/assignments/:id/submissions", s.handleGetAssignmentSubmissions, s.authorizedOne(model.ViewAnySubmissions))
	apiV1.PUT("/assignments/:id", s.handleUpdateAssignment, s.authorizedOne(model.UpdateAssignment))
	apiV1.DELETE("/assignments/:id", s.handleDeleteAssignment, s.authorizedOne(model.DeleteAssignment))

	apiV1.POST("/submissions", s.handleCreateSubmission, s.authorizedOne(model.CreateSubmission))
	apiV1.GET("/submissions/:id", s.handleGetSubmission, s.authorizedOne(model.ViewAnySubmissions))
	apiV1.PUT("/submissions", s.handleUpdateSubmission, s.authorizedOne(model.UpdateSubmission))
	apiV1.DELETE("/submissions/:id", s.handleDeleteSubmission, s.authorizedOne(model.DeleteSubmission))

	apiV1.POST("/media/upload", s.handleUploadMedia, s.authorizedOne(model.CreateMedia))

	grpHandlerName, grpcHandler := autogradv1connect.NewAutogradServiceHandler(
		service.NewService(s.gormDB),
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
