package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/usecase"
	"github.com/sirupsen/logrus"
)

// Server ..
type Server struct {
	exampleUsecase    usecase.ExampleUsecase
	userUsecase       usecase.UserUsecase
	assignmentUsecase usecase.AssignmentUsecase
	submissionUsecase usecase.SubmissionUsecase
	echo              *echo.Echo
	port              string
}

// NewServer ..
func NewServer(port string, opts ...Option) *Server {
	s := &Server{
		echo: echo.New(),
		port: port,
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

func (s *Server) routes() {
	s.echo.Static("/storage", "submission")
	s.echo.GET("/ping", s.handlePing)

	apiV1 := s.echo.Group("/api/v1")
	apiV1.POST("/users", s.handleCreateUser)
	apiV1.POST("/users/login", s.handleLogin)

	// example using auth middleware
	authorizeAdminStudent := []model.Role{model.RoleAdmin, model.RoleStudent}
	apiV1.GET("/example-private-data", s.handlePing, AuthMiddleware, s.authorizeByRoleMiddleware(authorizeAdminStudent))

	apiV1.POST("/assignments", s.handleCreateAssignment)
	apiV1.GET("/assignments/:ID", s.handleGetAssignment)
	apiV1.GET("/assignments", s.handleGetAssignments)
	apiV1.PUT("/assignments", s.handleUpdateAssignment)
	apiV1.DELETE("/assignments/:ID", s.handleDeleteAssignment)

	apiV1.POST("/submissions", s.handleCreateSubmission)
	apiV1.POST("/submissions/upload", s.handleUpload)
	apiV1.GET("/submissions/:assignmentID", s.handleGetAssignmentSubmission)
}

func (s *Server) handlePing(c echo.Context) error {
	s.exampleUsecase.Test()
	return c.String(http.StatusOK, "pong")
}
