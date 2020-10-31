package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/usecase"
)

// Server ..
type Server struct {
	exampleUsecase usecase.ExampleUsecase
	userUsecase    usecase.UserUsecase
	echo           *echo.Echo
	port           string
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
	s.echo.Start(":" + s.port)
}

func (s *Server) routes() {
	s.echo.GET("/ping", s.handlePing)

	apiV1 := s.echo.Group("/api/v1")
	apiV1.POST("/users", s.handleCreateUser)
	apiV1.POST("/users/login", s.handleLogin)

	// example using auth middleware
	authorizeAdminStudent := []model.Role{model.RoleAdmin, model.RoleStudent}
	apiV1.GET("/example-private-data", s.handlePing, AuthMiddleware, s.authorizeByRoleMiddleware(authorizeAdminStudent))
}

func (s *Server) handlePing(c echo.Context) error {
	s.exampleUsecase.Test()
	return c.String(http.StatusOK, "pong")
}
