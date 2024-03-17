package httpsvc

import (
	"context"
	"errors"
	"strings"

	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"github.com/fahmifan/autograd/pkg/service"
	"github.com/labstack/echo-contrib/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"
	_ "net/http/pprof"
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
	s.service.InternalCreateMacSandBoxRules()
	err := s.echo.Start(":" + s.port)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logs.Err(err, "Server: Start")
	}
}

// Stop server gracefully
func (s *Server) Stop(ctx context.Context) {
	if err := s.echo.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logs.ErrCtx(ctx, err, "Server: Stop")
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
	apiV1.GET("/rpc/activateManagedUser", s.handleActivateManagedUser)

	pprof.Register(s.echo, "/debug/pprof")

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
