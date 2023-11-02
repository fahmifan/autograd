package cpp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fahmifan/autograd/pkg/core/grading"
	"github.com/sirupsen/logrus"
)

var _ grading.Compiler = (*CPPCompiler)(nil)

type CPPCompiler struct {
}

func (c *CPPCompiler) Compile(inputPath string) (outPath string, err error) {
	outPath = path.Join(fmt.Sprintf("%s.bin", inputPath))
	args := strings.Split(fmt.Sprintf("%s -o %s", inputPath, outPath), " ")
	cmd := exec.Command("g++", args...)
	bt, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"output":  string(bt),
			"args":    args,
			"outPath": outPath,
		}).Error(err)
	}

	return
}

func (c *CPPCompiler) Run(bindPath string, input io.Reader, output io.Writer) (err error) {
	cmd := exec.Command(bindPath)

	buffErr := bytes.NewBuffer(nil)

	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = buffErr

	err = cmd.Run()
	if err != nil {
		return
	}

	if buffErr.Len() > 0 {
		err = fmt.Errorf(buffErr.String())
	}

	return
}

func (c *CPPCompiler) Remove(source string) error {
	err := os.Remove(source)
	if err != nil {
		logrus.Error(err)
	}

	return err
}
