package repository

import "github.com/sirupsen/logrus"

type SubmissionRepository interface {
	Test()
}

type submissionRepo struct {
}

func NewSubmissionRepo() SubmissionRepository {
	return &submissionRepo{}
}

func (e *submissionRepo) Test() {
	logrus.Warn("test: OK")
}
