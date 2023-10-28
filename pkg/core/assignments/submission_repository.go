package assignments

import (
	"context"
	"fmt"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionReader struct{}

func (SubmissionReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Submission, error) {
	subModel := dbmodel.Submission{}
	submitterModel := dbmodel.User{}
	sourceFileModel := dbmodel.File{}

	err := tx.WithContext(ctx).Where("id = ?", id).Take(&subModel).Error
	if err != nil {
		return Submission{}, fmt.Errorf("find submission: %w", err)
	}

	err = tx.WithContext(ctx).Where("id = ?", subModel.SubmittedBy).Take(&submitterModel).Error
	if err != nil {
		return Submission{}, fmt.Errorf("find submitter: %w", err)
	}

	assignment, err := AssignmentReader{}.FindByID(ctx, tx, subModel.AssignmentID)
	if err != nil {
		return Submission{}, fmt.Errorf("find assignment: %w", err)
	}

	err = tx.WithContext(ctx).Where("id = ?", subModel.FileID).Take(&sourceFileModel).Error
	if err != nil {
		return Submission{}, fmt.Errorf("find source file: %w", err)
	}

	return Submission{
		ID:         subModel.ID,
		Assignment: assignment,
		Submitter: Submitter{
			ID:     submitterModel.ID,
			Name:   submitterModel.Name,
			Active: submitterModel.Active == 1,
		},
		SourceFile: SubmissionFile{
			ID:  sourceFileModel.ID,
			URL: sourceFileModel.URL,
		},
	}, nil
}

type SubmissionWriter struct{}

func (SubmissionWriter) SaveNew(ctx context.Context, tx *gorm.DB, submission *Submission) error {
	model := dbmodel.Submission{
		Base: dbmodel.Base{
			ID:       submission.ID,
			Metadata: core.NewModelMetadata(submission.EntityMeta),
		},
		AssignmentID: submission.Assignment.ID,
		FileID:       submission.SourceFile.ID,
		SubmittedBy:  submission.Submitter.ID,
		Grade:        submission.Grade,
		Feedback:     submission.Feedback,
	}
	return tx.WithContext(ctx).Create(model).Error
}

func (SubmissionWriter) Save(ctx context.Context, tx *gorm.DB, submission *Submission) error {
	model := dbmodel.Submission{
		Base: dbmodel.Base{
			ID:       submission.ID,
			Metadata: core.NewModelMetadata(submission.EntityMeta),
		},
		AssignmentID: submission.Assignment.ID,
		FileID:       submission.SourceFile.ID,
		SubmittedBy:  submission.Submitter.ID,
		Grade:        submission.Grade,
		Feedback:     submission.Feedback,
	}
	return tx.WithContext(ctx).Save(model).Error
}

func (SubmissionWriter) Delete(ctx context.Context, tx *gorm.DB, submission *Submission) error {
	return tx.WithContext(ctx).UpdateColumn("deleted_at", submission.DeletedAt).Error
}
