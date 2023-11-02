package httpsvc

import (
	"context"
	"strings"

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
	echo            *echo.Echo
	gormDB          *gorm.DB
	port            string
	staticMediaPath string
	service         *service.Service
	jwtKey          string
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

	apiV1 := s.echo.Group("/api/v1")
	apiV1.POST("/rpc/saveMedia", s.handleSaveMedia)

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
