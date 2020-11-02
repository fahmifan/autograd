package grader

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// CompilerType types of compiler
type CompilerType int

// compiler types
const (
	CPPCompiler = CompilerType(0)
)

// Compiler ..
type Compiler interface {
	Compile(inputPath string) (outPath string, err error)
	Run(source, input string) (out string, err error)
	Remove(source string) error
}

// NewCompiler compiler factory
func NewCompiler(t CompilerType) Compiler {
	switch t {
	case CPPCompiler:
		return &cppCompiler{}
	default:
		return nil
	}
}

type cppCompiler struct {
}

// Compile compile programs
func (c *cppCompiler) Compile(inputPath string) (outPath string, err error) {
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

// Run the binary with input as arguments and return the output
func (c *cppCompiler) Run(source, input string) (out string, err error) {
	inputs := strings.Split(input, " ")
	input = strings.Join(inputs, "\n")
	cmd := exec.Command(source, inputs...)

	var buffOut bytes.Buffer
	var buffErr bytes.Buffer

	cmd.Stdin = bytes.NewBuffer([]byte(input))
	cmd.Stdout = &buffOut
	cmd.Stderr = &buffErr

	err = cmd.Run()
	if err != nil {
		return
	}

	out = strings.TrimSpace(buffOut.String())
	return
}

func (c *cppCompiler) Remove(source string) error {
	err := os.Remove(source)
	if err != nil {
		logrus.Error(err)
	}

	return err
}
