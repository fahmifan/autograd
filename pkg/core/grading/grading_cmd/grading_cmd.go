package grading_cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/grading"
	"github.com/fahmifan/autograd/pkg/core/grading/cpp"
	"github.com/google/uuid"
)

type GradingCmd struct {
	*core.Ctx
}

type InternalGradeSubmissionRequest struct {
	SubmissionID uuid.UUID
}

type InternalGradeSubmissionResult struct {
	SubmissionID uuid.UUID
}

// InternalCreateMacSandBoxRules creates mac sandbox rules
// and store it to local disk.
// This function is not concurrency safe.
func (cmd *GradingCmd) InternalCreateMacSandBoxRules() error {
	ruleFile, err := os.OpenFile(grading.RuleName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	rd := bufio.NewReader(ruleFile)
	line, _, err := rd.ReadLine()
	if err != nil && err != io.EOF {
		return fmt.Errorf("read line: %w", err)
	}

	if string(line) != "" {
		return nil
	}

	err = grading.CreateMacSandboxRules(ruleFile)
	if err != nil {
		return fmt.Errorf("create mac sandbox rules: %w", err)
	}

	return nil
}

func (cmd *GradingCmd) InternalGradeSubmission(
	ctx context.Context,
	req InternalGradeSubmissionRequest,
) (InternalGradeSubmissionResult, error) {
	submission, err := grading.SubmissionReader{}.FindByID(ctx, cmd.GormDB, cmd.ObjectStorer, cmd.RootDir, req.SubmissionID)
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("find submission: %w", err)
	}
	defer func() {
		submission.SubmissionFile.File.Close()
		submission.Assignment.CaseInputFile.File.Close()
		submission.Assignment.CaseOutputFile.File.Close()
	}()

	dir := path.Join(os.TempDir(), cmd.RootDir)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("create temp dir: %w", err)
	}

	submissionFilePath := path.Join(dir, submission.SubmissionFile.FileName)
	// store submission file to local disk.
	// use local scope to defer close the file
	{
		file, err := os.Create(submissionFilePath)
		if err != nil {
			return InternalGradeSubmissionResult{}, fmt.Errorf("create temp submission file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(file, submission.SubmissionFile.File)
		if err != nil {
			return InternalGradeSubmissionResult{}, fmt.Errorf("copy submission file: %w", err)
		}
	}

	compiler := &cpp.CPPCompiler{}

	gradeRes, err := grading.Grade(grading.GradeRequest{
		Compiler:       compiler,
		SourceCodePath: grading.SourceCodePath(submissionFilePath),
		Inputs:         submission.Assignment.CaseInputFile.File,
		Expecteds:      submission.Assignment.CaseOutputFile.File,
	})
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("grade: %w", err)
	}

	submission = submission.SaveGrade(time.Now(), gradeRes)

	err = grading.SubmissionWriter{}.Update(ctx, cmd.GormDB, &submission)
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("update submission: %w", err)
	}

	return InternalGradeSubmissionResult{
		SubmissionID: submission.ID,
	}, nil
}
