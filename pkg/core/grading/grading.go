package grading

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core/grading/podman"
	"github.com/google/uuid"
)

type Language string

const (
	LanguageCPP Language = "cpp"
)

type Compiler interface {
	Run(srcCodePath SourceCodePath, input io.Reader, output io.Writer) (err error)
}

type SourceCodePath string
type SourceCodeDir string
type RelativeFilename string

type GradeRequest struct {
	Compiler       Compiler
	SourceCodePath SourceCodePath
	Expecteds      io.Reader
	Inputs         io.Reader
	Submission     Submission
}

type GradeResult struct {
	Outputs  []string
	Corrects []bool
}

func Grade(arg GradeRequest) (GradeResult, error) {
	compiler := arg.Compiler
	result := GradeResult{}

	stdout := bytes.NewBuffer(nil)
	err := compiler.Run(arg.SourceCodePath, arg.Inputs, stdout)
	if err != nil {
		return GradeResult{}, fmt.Errorf("Grade: run: %w", err)
	}

	outputs := strings.Split(stdout.String(), "\n")
	expectedbuf, err := io.ReadAll(arg.Expecteds)
	if err != nil {
		return GradeResult{}, fmt.Errorf("Grade: read expecteds: %w", err)
	}

	expecteds := strings.Split(string(expectedbuf), "\n")
	if len(outputs) != len(expecteds) {
		return GradeResult{}, fmt.Errorf("Grade: expecteds and outputs length mismatch")
	}

	for i, output := range outputs {
		result.Outputs = append(result.Outputs, output)
		result.Corrects = append(result.Corrects, output == expecteds[i])
	}

	return result, nil
}

type GradeRequestV2 struct {
	Compiler         Compiler
	RelativeFilename RelativeFilename
	SourceCodeDir    SourceCodeDir
	Expecteds        io.Reader
	Inputs           io.Reader
	Submission       Submission
}

func GradeV2(arg GradeRequestV2) (GradeResult, error) {
	compiler := podman.CPP{
		MountDir:        string(arg.SourceCodeDir),
		ProgramFileName: string(arg.RelativeFilename),
		Input:           arg.Inputs,
	}

	result := GradeResult{}

	res, err := compiler.Run()
	if err != nil {
		return GradeResult{}, fmt.Errorf("Grade: run: %w", err)
	}

	outputs := strings.Split(string(res.Output()), "\n")
	expectedbuf, err := io.ReadAll(arg.Expecteds)
	if err != nil {
		return GradeResult{}, fmt.Errorf("Grade: read expecteds: %w", err)
	}

	expecteds := strings.Split(string(expectedbuf), "\n")
	if len(outputs) != len(expecteds) {
		return GradeResult{}, fmt.Errorf("Grade: expecteds and outputs length mismatch")
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
	Grade          int32
	Feedback       string
	UpdatedAt      time.Time
	IsGraded       bool
}

func (submission Submission) SaveGrade(now time.Time, grade GradeResult) Submission {
	var sum int32
	for _, correct := range grade.Corrects {
		if correct {
			sum++
		}
	}

	gradeScore := sum * 100 / int32(len(grade.Corrects))

	submission.Grade = gradeScore
	submission.UpdatedAt = now
	submission.IsGraded = true

	return submission
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
	FilePath string
	File     io.ReadCloser
}

type CaseInputFile struct {
	File io.ReadCloser
}

type CaseOutputFile struct {
	File io.ReadCloser
}
