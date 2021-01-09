package grader

// Option ..
type Option func(g *Grader)

// WithAssignmentUsecase ..
func WithAssignmentUsecase(a AssignmentUsecase) Option {
	return func(g *Grader) {
		g.assignmentUsecase = a
	}
}

// WithSubmissionUsecase ..
func WithSubmissionUsecase(s SubmissionUsecase) Option {
	return func(g *Grader) {
		g.submisisonUsecase = s
	}
}
