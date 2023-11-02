package grading

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/google/uuid"
)

type Language string

const (
	LanguageCPP Language = "cpp"
)

type Compiler interface {
	Compile(inputPath string) (outPath string, err error)
	Run(bindPath string, input io.Reader, output io.Writer) (err error)
	Remove(source string) error
}

type GradeRequest struct {
	Compiler       Compiler
	SourceCodePath string
	Expecteds      io.Reader
	Inputs         io.Reader
}

type GradeResult struct {
	Outputs  []string
	Corrects []bool
}

func Grade(arg GradeRequest) (GradeResult, error) {
	compiler := arg.Compiler

	binPath, err := compiler.Compile(arg.SourceCodePath)
	if err != nil {
		return GradeResult{}, fmt.Errorf("compile: %w", err)
	}

	defer func() {
		if err := compiler.Remove(binPath); err != nil {
			logs.Err(err, "path", binPath)
		}
	}()

	result := GradeResult{}

	stdout := bytes.NewBuffer(nil)
	err = compiler.Run(binPath, arg.Inputs, stdout)
	if err != nil {
		return GradeResult{}, fmt.Errorf("run: %w", err)
	}

	outputs := strings.Split(stdout.String(), "\n")
	expectedbuf, err := io.ReadAll(arg.Expecteds)
	if err != nil {
		return GradeResult{}, fmt.Errorf("read expecteds: %w", err)
	}

	expecteds := strings.Split(string(expectedbuf), "\n")
	if len(outputs) != len(expecteds) {
		return GradeResult{}, fmt.Errorf("expecteds and outputs length mismatch")
	}

	for i, output := range outputs {
		result.Outputs = append(result.Outputs, output)
		result.Corrects = append(result.Corrects, output == expecteds[i])
	}

	return result, nil
}

type Submission struct {
	ID             uuid.UUID
	Student        Student
	Assigner       Assigner
	Assignment     Assignment
	SubmissionFile SubmissionFile
	Grade          int64
	Feedback       string
}

type Assigner struct {
	ID   uuid.UUID
	Name string
}

type Student struct {
	ID     uuid.UUID
	Name   string
	Active bool
}

type Assignment struct {
	ID             uuid.UUID
	DeadlineAt     time.Time
	CaseInputFile  CaseInputFile
	CaseOutputFile CaseOutputFile
}

type SubmissionFile struct {
	FileName string
	File     io.ReadCloser
}

type CaseInputFile struct {
	File io.ReadCloser
}

type CaseOutputFile struct {
	File io.ReadCloser
}
