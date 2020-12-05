package httpsvc

import "github.com/miun173/autograd/usecase"

// Option ..
type Option func(*Server)

// WithExampleUsecase ..
func WithExampleUsecase(ex usecase.ExampleUsecase) Option {
	return func(s *Server) {
		s.exampleUsecase = ex
	}
}

// WithUserUsecase ..
func WithUserUsecase(u usecase.UserUsecase) Option {
	return func(s *Server) {
		s.userUsecase = u
	}
}

// WithSubmissionUsecase ..
func WithSubmissionUsecase(sub usecase.SubmissionUsecase) Option {
	return func(s *Server) {
		s.submissionUsecase = sub
	}
}
