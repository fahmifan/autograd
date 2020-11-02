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

func (g *graderImpl) Grade(source string, inputs, expecteds []string) (outputs []string, corrects []bool, err error) {
	outPath, err := g.compiler.Compile(source)
	if err != nil {
		logrus.WithField("source", source).Error(err)
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
			return nil, nil, err
		}

		outputs = append(outputs, out)
		correct := false
		if out == expected {
			correct = true
		}

		corrects = append(corrects, correct)
	}

	return
}

func (g *graderImpl) removeCompiled(path string) {
	if err := g.compiler.Remove(path); err != nil {
		logrus.WithField("path", path).Error(err)
	}
}
