package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Config ..
type Config struct {
	echo       *echo.Echo
	Port       int
	Debug      bool
	APIBaseURL string
}

// Server ..
type Server struct {
	*Config
}

// NewServer ..
func NewServer(cfg *Config) *Server {
	cfg.echo = echo.New()
	return &Server{cfg}
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

	// static
	s.echo.Static("/public", "web/public")
	s.echo.GET("/public/js/config.js", s.webConfig())

	s.echo.GET("", renderHTML("home.html")).Name = _routes.HomePage
	s.echo.GET("/auth/login", renderHTML("auth/login.html")).Name = _routes.LoginPage
}

func (s *Server) webConfig() echo.HandlerFunc {
	cfg := fmt.Sprintf(`
		export const config = {
			"API_URL": "%s",
		}`,
		s.APIBaseURL,
	)
	return func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/javascript", []byte(cfg))
	}
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
