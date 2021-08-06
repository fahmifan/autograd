package model

// WorkerBroker ..
type WorkerBroker interface {
	EnqueueJobGradeSubmission(submissionID string) error
}
