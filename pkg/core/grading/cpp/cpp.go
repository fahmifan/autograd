package cpp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/fahmifan/autograd/pkg/core/grading"
	"github.com/fahmifan/autograd/pkg/logs"
)

var _ grading.Compiler = (*CPPCompiler)(nil)

type CPPCompiler struct {
}

func (c *CPPCompiler) compile(inputPath grading.SourceCodePath) (outPath string, err error) {
	outPath = path.Join(fmt.Sprintf("%s.bin", inputPath))
	args := strings.Split(fmt.Sprintf("%s -o %s", inputPath, outPath), " ")
	cmd := exec.Command("g++", args...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("comile: %w", err)
	}

	return outPath, nil
}

func (c *CPPCompiler) Run(srcCodePath grading.SourceCodePath, input io.Reader, output io.Writer) (err error) {
	return c.run(srcCodePath, input, output)
}

func (c *CPPCompiler) run(srcCodePath grading.SourceCodePath, input io.Reader, output io.Writer) (err error) {
	binPath, err := c.compile(srcCodePath)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	defer func() {
		if err := c.remove(binPath); err != nil {
			logs.Err(err, "path", "binPath: ", string(binPath), "srcCodePath: ", string(srcCodePath))
		}
	}()

	cmd := exec.Command(binPath)

	if runtime.GOOS == "darwin" {
		sandboxRulePath := grading.RuleName
		cmd.Path = "/usr/bin/sandbox-exec"
		cmd = exec.Command("/usr/bin/sandbox-exec", "-f", sandboxRulePath, binPath)
	}

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

func (c *CPPCompiler) remove(source string) error {
	return os.Remove(source)
}
