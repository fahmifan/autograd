package assignments

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileReader struct{}

func (FileReader) FindCaseFiles(ctx context.Context, tx *gorm.DB, ids []uuid.UUID) ([]CaseFile, error) {
	var files []dbmodel.File
	err := tx.Table("files").Where("id IN (?)", ids).Find(&files).Error
	if err != nil {
		return nil, err
	}

	return toCaseFiles(files), nil
}

func toCaseFiles(files []dbmodel.File) []CaseFile {
	var caseFiles []CaseFile
	for _, file := range files {
		caseFiles = append(caseFiles, toCaseFile(file))
	}

	return caseFiles
}
