package assignments

import (
	"context"
	"errors"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type AssignmentWriter struct{}

func (AssignmentWriter) Save(ctx context.Context, tx *gorm.DB, assignment Assignment) error {
	model := dbmodel.Assignment{
		Base: dbmodel.Base{
			ID: assignment.ID,
		},
		AssignedBy:       assignment.Assigner.ID,
		Name:             assignment.Name,
		Description:      assignment.Description,
		CaseInputFileID:  assignment.CaseInputFile.ID,
		CaseOutputFileID: assignment.CaseOutputFile.ID,
	}

	return tx.Table("assignments").Create(&model).Error
}

type AssignmentReader struct{}

func (AssignmentReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Assignment, error) {
	assignment := dbmodel.Assignment{}
	err := tx.Table("assignments").Where("id = ?", id).Take(&assignment).Error
	if err != nil {
		return Assignment{}, err
	}

	user := dbmodel.User{}
	err = tx.Table("users").Where("id = ?", id).Take(&user).Error
	if err != nil {
		return Assignment{}, err
	}

	files := []dbmodel.File{}
	fileIDs := []uuid.UUID{assignment.CaseInputFileID, assignment.CaseOutputFileID}
	err = tx.Table("files").Where("id IN (?)", fileIDs).Find(&files).Error
	if err != nil {
		return Assignment{}, err
	}

	if len(files) != 2 {
		return Assignment{}, err
	}

	caseInputFile, _, found := lo.FindIndexOf(files, func(file dbmodel.File) bool {
		return file.Type == dbmodel.FileTypeAssignmentCaseInput
	})
	if !found {
		return Assignment{}, errors.New("case input file not found")
	}

	caseOutputFile, _, found := lo.FindIndexOf(files, func(file dbmodel.File) bool {
		return file.Type == dbmodel.FileTypeAssignmentCaseOutput
	})
	if !found {
		return Assignment{}, errors.New("case output file not found")
	}

	return toAssignment(assignment, user, caseInputFile, caseOutputFile), err
}

type AssignerReader struct{}

func (AssignerReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Assigner, error) {
	user := dbmodel.User{}
	err := tx.Table("users").Where("id = ?", id).Take(&user).Error
	return toAssigner(user), err
}

func toAssigner(user dbmodel.User) Assigner {
	return Assigner{
		ID:     user.ID,
		Name:   user.Name,
		Active: user.Active == 1,
	}
}

func toCaseFile(file dbmodel.File) CaseFile {
	return CaseFile{
		ID:   file.ID,
		URL:  file.URL,
		Type: file.Type,
	}
}

func toAssignment(
	model dbmodel.Assignment,
	user dbmodel.User,
	inputFile dbmodel.File,
	outputFile dbmodel.File,
) Assignment {
	return Assignment{
		ID:             model.ID,
		Name:           model.Name,
		Description:    model.Description,
		Assigner:       toAssigner(user),
		CaseInputFile:  toCaseFile(inputFile),
		CaseOutputFile: toCaseFile(outputFile),
		EntityMeta:     toEntityMeta(model.Base),
	}
}

func toEntityMeta(base dbmodel.Base) core.EntityMeta {
	return core.EntityMeta{
		CreatedAt: base.CreatedAt,
		UpdatedAt: base.UpdatedAt,
	}
}
