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
	"github.com/fahmifan/autograd/pkg/core/grading/podman"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	res := InternalGradeSubmissionResult{}
	err := core.Transaction(ctx, cmd.Ctx, func(tx *gorm.DB) (err error) {
		res, err = cmd.InternalGradeSubmissionTx(ctx, tx, req)
		if err != nil {
			return fmt.Errorf("InternalGradeSubmission: Transaction: InternalGradeSubmissionTx: %w", err)
		}
		return nil
	})
	return res, err
}

func (cmd *GradingCmd) InternalGradeSubmissionTx(
	ctx context.Context,
	tx *gorm.DB,
	req InternalGradeSubmissionRequest,
) (InternalGradeSubmissionResult, error) {
	submission, err := grading.SubmissionReader{}.FindByID(ctx, cmd.GormDB, cmd.ObjectStorer, cmd.RootDir, req.SubmissionID)
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("InternalGradeSubmissionTx: find submission: %w", err)
	}
	defer func() {
		submission.SubmissionFile.File.Close()
		submission.Assignment.CaseInputFile.File.Close()
		submission.Assignment.CaseOutputFile.File.Close()
	}()

	submissionFilePath := submission.SubmissionFile.FilePath
	fileDir, _ := path.Split(path.Join(cmd.RootDir, submissionFilePath))

	compiler := &podman.CPP{}

	gradeRes, err := grading.Grade(grading.GradeRequest{
		Compiler:         compiler,
		RelativeFilename: grading.RelativeFilename(submission.SubmissionFile.FileName),
		SourceCodeDir:    grading.SourceCodeDir(fileDir),
		Inputs:           submission.Assignment.CaseInputFile.File,
		Expecteds:        submission.Assignment.CaseOutputFile.File,
	})
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("InternalGradeSubmissionTx: grade: %w", err)
	}

	submission = submission.SaveGrade(time.Now(), gradeRes)

	err = grading.SubmissionWriter{}.Update(ctx, cmd.GormDB, &submission)
	if err != nil {
		return InternalGradeSubmissionResult{}, fmt.Errorf("InternalGradeSubmissionTx: update submission: %w", err)
	}

	return InternalGradeSubmissionResult{
		SubmissionID: submission.ID,
	}, nil
}
