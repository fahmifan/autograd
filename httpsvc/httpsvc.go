package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/usecase"
	usecaseIface "github.com/miun173/autograd/usecase/iface"
	"github.com/sirupsen/logrus"
)

// Server ..
type Server struct {
	echo              *echo.Echo
	port              string
	staticMediaPath   string
	exampleUsecase    usecase.ExampleUsecase
	userUsecase       usecase.UserUsecase
	assignmentUsecase usecase.AssignmentUsecase
	submissionUsecase usecase.SubmissionUsecase
	mediaUsecase      usecaseIface.MediaUsecase
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

func (s *Server) routes() {
	s.echo.GET("/ping", s.handlePing)

	// TODO: add auth for private static
	s.echo.Static("/storage", "submission")
	s.echo.Static("/media", s.staticMediaPath)

	apiV1 := s.echo.Group("/api/v1")
	apiV1.POST("/users", s.handleCreateUser)
	apiV1.POST("/users/login", s.handleLogin)

	// example using auth middleware
	authorizeAdminStudent := []model.Role{model.RoleAdmin, model.RoleStudent}
	apiV1.GET("/example-private-data", s.handlePing, AuthMiddleware, s.authorizeByRoleMiddleware(authorizeAdminStudent))

	apiV1.POST("/assignments", s.handleCreateAssignment)
	apiV1.GET("/assignments", s.handleGetAssignments)
	apiV1.GET("/assignments/:ID", s.handleGetAssignment)
	apiV1.GET("/assignments/:ID/submissions", s.handleGetAssignmentSubmissions)
	apiV1.PUT("/assignments", s.handleUpdateAssignment)
	apiV1.DELETE("/assignments/:ID", s.handleDeleteAssignment)

	apiV1.POST("/submissions", s.handleCreateSubmission)
	apiV1.POST("/submissions/upload", s.handleUpload)
	apiV1.GET("/submissions/:ID", s.handleGetSubmission)
	apiV1.PUT("/submissions", s.handleUpdateSubmission)
	apiV1.DELETE("/submissions/:ID", s.handleDeleteSubmission)

	apiV1.POST("/media/upload", s.handleUploadMedia)
}

func (s *Server) handlePing(c echo.Context) error {
	s.exampleUsecase.Test()
	return c.String(http.StatusOK, "pong")
}
