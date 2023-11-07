package httpsvc

import (
	"context"
	"log"
	"strings"

	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"github.com/fahmifan/autograd/pkg/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server ..
type Server struct {
	echo    *echo.Echo
	port    string
	service *service.Service
	jwtKey  string
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
	log.Fatal(s.echo.Start(":" + s.port))
}

// Stop server gracefully
func (s *Server) Stop(ctx context.Context) {
	if err := s.echo.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) routes() {
	s.echo.Use(
		middleware.CORS(),
		s.addUserToCtx,
		logs.EchoRequestID(),
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
