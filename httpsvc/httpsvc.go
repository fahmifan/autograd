package httpsvc

import (
	"context"
	"net/http"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Server ..
type Server struct {
	echo              *echo.Echo
	port              string
	staticMediaPath   string
	userUsecase       usecase.UserUsecase
	assignmentUsecase usecase.AssignmentUsecase
	submissionUsecase usecase.SubmissionUsecase
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

	authAdminStudent := s.authorizeByRoleMiddleware([]model.Role{model.RoleAdmin, model.RoleStudent})

	apiV1.POST("/assignments", s.handleCreateAssignment)
	apiV1.GET("/assignments", s.handleGetAssignments, authAdminStudent)
	apiV1.GET("/assignments/:ID", s.handleGetAssignment, authAdminStudent)
	apiV1.GET("/assignments/:ID/submissions", s.handleGetAssignmentSubmissions)
	apiV1.PUT("/assignments", s.handleUpdateAssignment)
	apiV1.DELETE("/assignments/:ID", s.handleDeleteAssignment)

	apiV1.POST("/submissions", s.handleCreateSubmission, authAdminStudent)
	apiV1.GET("/submissions/:ID", s.handleGetSubmission, authAdminStudent)
	apiV1.PUT("/submissions", s.handleUpdateSubmission)
	apiV1.DELETE("/submissions/:ID", s.handleDeleteSubmission)

	apiV1.POST("/media/upload", s.handleUploadMedia, authAdminStudent)
}

func (s *Server) handlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
