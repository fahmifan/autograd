package repository

import "github.com/sirupsen/logrus"

// ExampleRepository ..
type ExampleRepository interface {
	Test()
}

type exampleRepo struct {
}

// NewExampleRepo ..
func NewExampleRepo() ExampleRepository {
	return &exampleRepo{}
}

func (e *exampleRepo) Test() {
	logrus.Warn("test: OK")
}
