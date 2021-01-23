package worker

import (
	"github.com/fahmifan/autograd/model"
)

// Option ..
type Option func(*Worker)

// WithGrader ..
func WithGrader(gr model.GraderUsecase) Option {
	return func(c *Worker) {
		c.grader = gr
	}
}

// WithSubmissionUsecase ..
func WithSubmissionUsecase(s model.SubmissionUsecase) Option {
	return func(c *Worker) {
		c.submission = s
	}
}

// WithAssignmentUsecase ..
func WithAssignmentUsecase(a model.AssignmentUsecase) Option {
	return func(c *Worker) {
		c.assignment = a
	}
}
