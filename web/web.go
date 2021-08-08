package web

import (
	"context"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Echo *echo.Echo
	Port string
}

func (s *Server) Run() error {
	s.route()
	logrus.Info(s.Port)
	return s.Echo.Start(":" + s.Port)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

var _routes = struct {
	HomePage  string
	LoginPage string
}{
	HomePage:  "HomePage",
	LoginPage: "LoginPage",
}

func (s *Server) route() {
	s.Echo.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("web/templates/**/*.html")),
	}

	s.Echo.Static("/public", "web/public")
	s.Echo.GET("", renderHTML("home.html")).Name = _routes.HomePage
	s.Echo.GET("/auth/login", renderHTML("login.html")).Name = _routes.LoginPage
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
