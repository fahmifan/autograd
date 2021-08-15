package httpsvc

import (
	"github.com/fahmifan/autograd/model"
)

// Option ..
type Option func(*Server)

// WithUserUsecase ..
func WithUserUsecase(u model.UserUsecase) Option {
	return func(s *Server) {
		s.userUsecase = u
	}
}

// WithAssignmentUsecase ..
func WithAssignmentUsecase(a model.AssignmentUsecase) Option {
	return func(s *Server) {
		s.assignmentUsecase = a
	}
}

// WithSubmissionUsecase ..
func WithSubmissionUsecase(sub model.SubmissionUsecase) Option {
	return func(s *Server) {
		s.submissionUsecase = sub
	}
}

// WithObjectStorer ..
func WithObjectStorer(obs model.ObjectStorer) Option {
	return func(s *Server) {
		s.objectStorer = obs
	}
}

// WithSessionRepository ..
func WithSessionRepository(sr model.SessionRespository) Option {
	return func(s *Server) {
		s.sessionRepo = sr
	}
}
