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
		return "", fmt.Errorf("comile: %w", err)
	}

	return string(bt), nil
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
	return os.Remove(source)
}
