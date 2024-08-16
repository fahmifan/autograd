package grading

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Language string

const (
	LanguageCPP Language = "cpp"
)

type Second int

func (s Second) String() string {
	return fmt.Sprintf("%ds", s)
}

type Mib int

func (mib Mib) String() string {
	return fmt.Sprintf("%dm", mib)
}

type RunnerArg struct {
	Input           io.Reader
	MountDir        string
	ProgramFileName string
	MemLimit        Mib
	RunTimeout      Second
}

type Runner interface {
	Run(arg RunnerArg) (RunResult, error)
}

type RunResult struct {
	Output []byte
}

type SourceCodePath string
type SourceCodeDir string
type RelativeFilename string

type GradeResult struct {
	Outputs  []string
	Corrects []bool
}

type GradeRequest struct {
	Compiler         Runner
	RelativeFilename RelativeFilename
	SourceCodeDir    SourceCodeDir
	Expecteds        io.Reader
	Inputs           io.Reader
	Submission       Submission
}

func Grade(arg GradeRequest) (GradeResult, error) {
	compiler := arg.Compiler

	runRes, err := compiler.Run(RunnerArg{
		MountDir:        string(arg.SourceCodeDir),
		ProgramFileName: string(arg.RelativeFilename),
		Input:           arg.Inputs,
		MemLimit:        100,
		RunTimeout:      10,
	})
	if err != nil {
		return GradeResult{}, fmt.Errorf("grade: run: %w", err)
	}

	outputs := strings.Split(string(runRes.Output), "\n")
	expectedbuf, err := io.ReadAll(arg.Expecteds)
	if err != nil {
		return GradeResult{}, fmt.Errorf("grade: read expecteds: %w", err)
	}

	expecteds := strings.Split(string(expectedbuf), "\n")
	if len(outputs) != len(expecteds) {
		return GradeResult{}, fmt.Errorf("grade: expecteds and outputs length mismatch")
	}

	result := GradeResult{}
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
