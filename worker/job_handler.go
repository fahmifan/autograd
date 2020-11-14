package worker

import (
	"github.com/gocraft/work"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// job names
const (
	jobRunCode string = "run_code"
)

type jobHandler struct {
	*Config
}

func (*jobHandler) handleRunCode(job *work.Job) error {
	logrus.Warn("handleRunCode >>>", utils.Dump(job))
	return nil
}
