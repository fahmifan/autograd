package httpsvc

import (
	"context"
	"net/http"

	"github.com/fahmifan/autograd/model"

	_ "github.com/fahmifan/autograd/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Server ..
type Server struct {
	echo              *echo.Echo
	port              string
	staticMediaPath   string
	userUsecase       model.UserUsecase
	assignmentUsecase model.AssignmentUsecase
	submissionUsecase model.SubmissionUsecase
	objectStorer      model.ObjectStorer
	sessionRepo       model.SessionRespository
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
	err := s.echo.Start(":" + s.port)
	if err != nil && err != http.ErrServerClosed {
		logrus.Error(err)
		return
	}
	logrus.Info("api server stopped gracefully")
}

// Stop server gracefully
func (s *Server) Stop(ctx context.Context) {
	if err := s.echo.Shutdown(ctx); err != nil {
		logrus.Error(err)
	}
}

func (s *Server) routes() {
	s.echo.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	s.echo.GET("/ping", s.handlePing)
	s.echo.GET("/docs/swagger/*", echoSwagger.WrapHandler)

	apiV1 := s.echo.Group("/api/v1")

	apiV1.POST("/auth/login", s.handleLogin)
	apiV1.POST("/auth/refresh", s.handleRefreshToken)

	apiV1.POST("/users", s.handleCreateUser)

	apiV1.POST("/assignments", s.handleCreateAssignment, s.authorizedAny(model.CreateAssignment))
	apiV1.GET("/assignments", s.handleGetAllAssignments, s.authorizedAny(model.ViewAnyAssignments, model.ViewAssignment))
	apiV1.GET("/assignments/:id", s.handleGetAssignment, s.authorizedAny(model.ViewAssignment, model.ViewAnyAssignments))
	apiV1.GET("/assignments/:id/submissions", s.handleGetAssignmentSubmissions, s.authorizedAny(model.ViewAnySubmissions, model.ViewSubmission))
	apiV1.PUT("/assignments/:id", s.handleUpdateAssignment, s.authorizedAny(model.UpdateAssignment))
	apiV1.DELETE("/assignments/:id", s.handleDeleteAssignment, s.authorizedAny(model.DeleteAssignment))

	apiV1.POST("/submissions", s.handleCreateSubmission, s.authorizedAny(model.CreateSubmission))
	apiV1.GET("/submissions/:id", s.handleGetSubmission, s.authorizedAny(model.ViewSubmission, model.ViewAnySubmissions))
	apiV1.PUT("/submissions", s.handleUpdateSubmission, s.authorizedAny(model.UpdateSubmission))
	apiV1.DELETE("/submissions/:id", s.handleDeleteSubmission, s.authorizedAny(model.DeleteSubmission))

	apiV1.POST("/media", s.handleUploadMedia, s.authorizedAny(model.CreateMedia))
	apiV1.GET("/media/:filename", s.handleGetMedia)
}

func (s *Server) handlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
