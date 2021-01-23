package grader

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/miun173/autograd/model"
	"github.com/sirupsen/logrus"
)

// make sure implement interface
var _ model.GraderEngine = (*CPPGrader)(nil)

// CPPGrader implements model.Grader
type CPPGrader struct {
}

// Grade ..
func (c *CPPGrader) Grade(arg *model.GradingArg) (*model.GradingResult, error) {
	binPath, err := c.Compile(arg.SourceCodePath)
	if err != nil {
		logrus.WithField("source", arg.SourceCodePath).Error(err)
		return nil, err
	}

	defer func() {
		if err := c.Remove(binPath); err != nil {
			logrus.WithField("path", binPath).Error(err)
		}
	}()

	if len(arg.Inputs) != len(arg.Expecteds) {
		return nil, fmt.Errorf("expecteds & inputs not match in length")
	}

	result := &model.GradingResult{}
	for i := range arg.Inputs {
		input := arg.Inputs[i]
		expected := arg.Expecteds[i]

		out, err := c.Run(binPath, input)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result.Outputs = append(result.Outputs, out)
		result.Corrects = append(result.Corrects, out == expected)
	}

	return result, nil
}

// Compile compile programs
func (c *CPPGrader) Compile(inputPath string) (outPath string, err error) {
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
func (c *CPPGrader) Run(source, input string) (out string, err error) {
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

// Remove ..
func (c *CPPGrader) Remove(source string) error {
	err := os.Remove(source)
	if err != nil {
		logrus.Error(err)
	}

	return err
}
