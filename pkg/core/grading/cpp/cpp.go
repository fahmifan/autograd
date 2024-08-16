package cpp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/fahmifan/autograd/pkg/core/grading"
	"github.com/fahmifan/autograd/pkg/logs"
)

type CPPCompiler struct {
}

func (c *CPPCompiler) compile(inputPath grading.SourceCodePath) (outPath string, err error) {
	outPath = path.Join(fmt.Sprintf("%s.bin", inputPath))
	cmd := exec.Command("g++", string(inputPath), "-o", outPath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compile: %w", err)
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
			logs.Err(err, "path", "binPath: ", string(binPath))
		}
		if err := c.remove(string(srcCodePath)); err != nil {
			logs.Err(err, "path", "srcCodePath: ", string(srcCodePath))
		}
	}()

	cmd := exec.Command(binPath)

	switch runtime.GOOS {
	case "darwin":
		sandboxRulePath := grading.RuleFilePath()
		cmd = exec.Command("/usr/bin/sandbox-exec", "-f", sandboxRulePath, binPath)
	case "linux":
		//
	}

	buffErr := bytes.NewBuffer(nil)

	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = buffErr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("CPPCompiler: run: %w", err)
	}

	if buffErr.Len() > 0 {
		err = fmt.Errorf(buffErr.String())
		return fmt.Errorf("CPPCompiler: stderr: %w", err)
	}

	return
}

func (c *CPPCompiler) remove(source string) error {
	return os.Remove(source)
}
