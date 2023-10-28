package httpsvc

import (
	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/service"
	"gorm.io/gorm"
)

type Option func(*Server)

func WithUserUsecase(u model.UserUsecase) Option {
	return func(s *Server) {
		s.userUsecase = u
	}
}

func WithAssignmentUsecase(a model.AssignmentUsecase) Option {
	return func(s *Server) {
		s.assignmentUsecase = a
	}
}

func WithSubmissionUsecase(sub model.SubmissionUsecase) Option {
	return func(s *Server) {
		s.submissionUsecase = sub
	}
}

func WithMediaUsecase(med model.MediaUsecase) Option {
	return func(s *Server) {
		s.mediaUsecase = med
	}
}

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
