package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Server ..
type Server struct {
	echo  *echo.Echo
	Port  int
	Debug bool
}

// NewServer ..
func NewServer(port int, debug bool) *Server {
	return &Server{
		echo:  echo.New(),
		Port:  port,
		Debug: debug,
	}
}

// Run() ..
func (s *Server) Run() error {
	s.route()
	logrus.Info(s.Port)
	return s.echo.Start(fmt.Sprintf(":%d", s.Port))
}

// Stop ..
func (s *Server) Stop(ctx context.Context) {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		logrus.Error(err)
	}
}

var _routes = struct {
	HomePage  string
	LoginPage string
}{
	HomePage:  "HomePage",
	LoginPage: "LoginPage",
}

func (s *Server) route() {
	s.echo.Renderer = NewRenderer("web/templates", s.Debug)

	s.echo.Static("/public", "web/public")
	s.echo.GET("", renderHTML("home.html")).Name = _routes.HomePage
	s.echo.GET("/auth/login", renderHTML("auth/login.html")).Name = _routes.LoginPage
}

func renderHTML(file string) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := c.Render(http.StatusOK, file, echo.Map{})
		if err != nil {
			logrus.Error(err)
		}
		return err
	}
}
