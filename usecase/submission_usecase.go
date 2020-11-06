package usecase

import (
	"github.com/miun173/autograd/repository"
	"github.com/sirupsen/logrus"
)

type SubmissionUsecase interface {
	Test()
}

type submissionUsecase struct {
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionUsecase(submissionRepo repository.SubmissionRepository) SubmissionUsecase {
	return &submissionUsecase{
		submissionRepo: submissionRepo,
	}
}

func (e *submissionUsecase) Test() {
	e.submissionRepo.Test()
	logrus.Warn("test: OK")
}
