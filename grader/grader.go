package grader

import (
	"github.com/sirupsen/logrus"
)

// Grader ..
type Grader interface {
	Grade(source string, inputs, expecteds []string) (outputs []string, corrects []bool, err error)
}

type graderImpl struct {
	compiler Compiler
}

// NewGrader ..
func NewGrader(c Compiler) Grader {
	return &graderImpl{
		compiler: c,
	}
}

func (g *graderImpl) Grade(sourceCode string, inputs, expecteds []string) (outputs []string, corrects []bool, err error) {
	outPath, err := g.compiler.Compile(sourceCode)
	if err != nil {
		logrus.WithField("source", sourceCode).Error(err)
		return
	}

	defer g.removeCompiled(outPath)

	if len(inputs) != len(expecteds) {
		return
	}

	for i := range inputs {
		input := inputs[i]
		expected := expecteds[i]

		out, err := g.compiler.Run(outPath, input)
		if err != nil {
			logrus.Error(err)
			return outputs, corrects, err
		}

		outputs = append(outputs, out)
		corrects = append(corrects, out == expected)
	}

	return
}

func (g *graderImpl) removeCompiled(path string) {
	if err := g.compiler.Remove(path); err != nil {
		logrus.WithField("path", path).Error(err)
	}
}
