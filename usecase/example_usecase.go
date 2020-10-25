package usecase

import (
	"github.com/miun173/autograd/repository"
	"github.com/sirupsen/logrus"
)

// ExampleUsecase ..
type ExampleUsecase interface {
	Test()
}

type exampleUsecase struct {
	exampleRepo repository.ExampleRepository
}

// NewExampleUsecase ..
func NewExampleUsecase(exampleRepo repository.ExampleRepository) ExampleUsecase {
	return &exampleUsecase{
		exampleRepo: exampleRepo,
	}
}

func (e *exampleUsecase) Test() {
	e.exampleRepo.Test()
	logrus.Warn("test: OK")
}
