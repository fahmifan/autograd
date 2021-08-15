package model

// Worker is a worker that can be invoked by a Worker Manager.
// Every methods in Worker interface should be defined in Broker too
type Worker interface {
	GradeSubmission(submissionID string) error
	GradeAssignment(assignmentID string) error
}

// Broker  can be used to queuing jobs to a Worker Manager
//
// A Worker Manager (WM) is a program that managed job queues and the workers
// e.g. github.com/gocraft/work or github.com/RichardKnop/machinery
//
// Each WM has different ways of calling the workers.
// So, to make it more flexible, we seperate the actual implementation of Worker with Broker & WM.
// When implements the Broker & WM it should provide adapters that call the actual Worker
type Broker interface {
	GradeSubmission(submissionID string) error
	GradeAssignment(assignmentID string) error
}
