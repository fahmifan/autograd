package model

// WorkerBroker ..
type WorkerBroker interface {
	EnqueueJobGradeSubmission(submissionID int64) error
}
