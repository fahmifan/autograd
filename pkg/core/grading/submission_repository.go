package grading

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type SubmissionReader struct{}

func (SubmissionReader) FindByID(ctx context.Context, tx *gorm.DB, objStorer core.ObjectStorer, rootDir string, id uuid.UUID) (Submission, error) {
	tx = tx.WithContext(ctx)

	submModel := dbmodel.Submission{}
	if err := tx.Where("id = ?", id).Take(&submModel).Error; err != nil {
		return Submission{}, fmt.Errorf("find submission: %w", err)
	}

	submFile := dbmodel.File{}
	if err := tx.Where("id = ?", submModel.FileID).Take(&submFile).Error; err != nil {
		return Submission{}, fmt.Errorf("find submission file: %w", err)
	}

	assignmentModel := dbmodel.Assignment{}
	if err := tx.Where("id = ?", submModel.AssignmentID).Take(&assignmentModel).Error; err != nil {
		return Submission{}, fmt.Errorf("find assignment: %w", err)
	}

	var assignmentFiles []dbmodel.File
	caseFileIDs := []uuid.UUID{assignmentModel.CaseInputFileID, assignmentModel.CaseOutputFileID}

	if err := tx.Where("id in (?)", caseFileIDs).Find(&assignmentFiles).Error; err != nil {
		return Submission{}, fmt.Errorf("find case files: %w", err)
	}

	studentModel := dbmodel.User{}
	if err := tx.Where("id = ?", submModel.SubmittedBy).Take(&studentModel).Error; err != nil {
		return Submission{}, fmt.Errorf("find student: %w", err)
	}

	assignerModel := dbmodel.User{}
	if err := tx.Where("id = ?", assignmentModel.AssignedBy).Take(&assignerModel).Error; err != nil {
		return Submission{}, fmt.Errorf("find assigner: %w", err)
	}

	caseInputModel, ok := lo.Find(assignmentFiles, func(file dbmodel.File) bool {
		return file.Type == dbmodel.FileTypeAssignmentCaseInput
	})
	if !ok {
		return Submission{}, fmt.Errorf("case input file not found")
	}

	caseOutputModel, ok := lo.Find(assignmentFiles, func(file dbmodel.File) bool {
		return file.Type == dbmodel.FileTypeAssignmentCaseOutput
	})
	if !ok {
		return Submission{}, fmt.Errorf("case input file not found")
	}

	submission := Submission{
		ID:        id,
		Grade:     submModel.Grade,
		Feedback:  submModel.Feedback,
		UpdatedAt: submModel.UpdatedAt.Time,
		SubmissionFile: SubmissionFile{
			FileName: submFile.Name,
			FilePath: submFile.Path,
		},
		Student: Student{
			ID:     studentModel.ID,
			Name:   studentModel.Name,
			Active: studentModel.IsActive(),
		},
		Assigner: Assigner{
			ID:   assignerModel.ID,
			Name: assignerModel.Name,
		},
		Assignment: Assignment{
			ID:         assignmentModel.ID,
			DeadlineAt: assignmentModel.DeadlineAt,
		},
	}

	// defer close file when error, so we don't leak it
	{
		var err error

		submissionFile, err := objStorer.Seek(ctx, path.Join(rootDir, submFile.Path))
		defer closeWhenErr(err, submissionFile)
		if err != nil {
			return Submission{}, fmt.Errorf("seek submission file: %w", err)
		}

		caseInputFile, err := objStorer.Seek(ctx, path.Join(rootDir, caseInputModel.Path))
		defer closeWhenErr(err, caseInputFile)
		if err != nil {
			return Submission{}, fmt.Errorf("seek case input: %w", err)
		}

		caseOutputFile, err := objStorer.Seek(ctx, path.Join(rootDir, caseOutputModel.Path))
		defer closeWhenErr(err, caseOutputFile)
		if err != nil {
			return Submission{}, fmt.Errorf("seek case output: %w", err)
		}

		submission.SubmissionFile.File = submissionFile
		submission.Assignment.CaseInputFile.File = caseInputFile
		submission.Assignment.CaseOutputFile.File = caseOutputFile
	}

	return submission, nil
}

type SubmissionWriter struct{}

func (SubmissionWriter) Update(ctx context.Context, tx *gorm.DB, submission *Submission) error {
	return tx.WithContext(ctx).Model(&dbmodel.Submission{}).
		Where("id = ?", submission.ID).
		UpdateColumns(map[string]any{
			"grade":      submission.Grade,
			"updated_at": submission.UpdatedAt,
			"is_graded":  intBool(submission.IsGraded),
		}).Error
}

func intBool(b bool) int {
	if b {
		return 1
	}
	return 0
}

func closeWhenErr(err error, readCloser io.ReadCloser) {
	if err != nil && readCloser != nil {
		readCloser.Close()
	}
}
