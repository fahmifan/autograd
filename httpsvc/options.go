package httpsvc

import (
	"github.com/fahmifan/autograd/pkg/service"
	"gorm.io/gorm"
)

type Option func(*Server)

func WithService(s *service.Service) Option {
	return func(srv *Server) {
		srv.service = s
	}
}

func WithGormDB(db *gorm.DB) Option {
	return func(s *Server) {
		s.gormDB = db
	}
}

func WithJWTKey(key string) Option {
	return func(s *Server) {
		s.jwtKey = key
	}
}
