package assignments

import (
	"context"
	"errors"
	"fmt"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbconn"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/xsqlc"
	"github.com/fahmifan/ulids"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type AssignmentWriter struct{}

func (AssignmentWriter) Save(ctx context.Context, tx *gorm.DB, assignment Assignment) error {
	query, err := dbconn.NewSqlcFromGorm(tx)
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "NewSqlcFromGorm")
	}

	err = query.SaveAssignment(ctx, xsqlc.SaveAssignmentParams{
		ID:               assignment.ID.String(),
		AssignedBy:       assignment.Assigner.ID.String(),
		Name:             assignment.Name,
		Description:      assignment.Description,
		CaseInputFileID:  assignment.CaseInputFile.ID.String(),
		CaseOutputFileID: assignment.CaseOutputFile.ID.String(),
		Template:         assignment.Template,
		DeadlineAt:       assignment.DeadlineAt,
		CreatedAt:        assignment.CreatedAt,
		UpdatedAt:        assignment.UpdatedAt,
		DeletedAt:        assignment.DeletedAt.NullTime,
	})
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "AssignmentWriter: SaveAssignment")
	}

	err = query.SaveAssignmentCourse(ctx, xsqlc.SaveAssignmentCourseParams{
		CourseID:     assignment.Course.ID.String(),
		AssignmentID: assignment.ID.String(),
	})
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "AssignmentWriter: SaveAssignmentCourse")
	}

	return nil
}

func (wr AssignmentWriter) Update(ctx context.Context, tx *gorm.DB, assignment Assignment) error {
	return wr.Save(ctx, tx, assignment)
}

type AssignmentReader struct{}

func (AssignmentReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Assignment, error) {
	query, err := dbconn.NewSqlcFromGorm(tx)
	if err != nil {
		return Assignment{}, logs.ErrWrapCtx(ctx, err, "NewSqlcFromGorm")
	}

	assignment, err := query.FindAssignmentByID(ctx, id.String())
	if err != nil {
		return Assignment{}, logs.ErrWrapCtx(ctx, err, "find assignment")
	}

	course, err := query.FindCourseByID(ctx, assignment.CourseID)
	if err != nil {
		return Assignment{}, logs.ErrWrapCtx(ctx, err, "find course")
	}

	user := dbmodel.User{}
	err = tx.Table("users").Where("id = ?", assignment.AssignedBy).Take(&user).Error
	if err != nil {
		return Assignment{}, logs.ErrWrapCtx(ctx, err, "find user")
	}

	files := []dbmodel.File{}
	fileIDs := []string{assignment.CaseInputFileID, assignment.CaseOutputFileID}
	err = tx.Table("files").Where("id IN (?)", fileIDs).Find(&files).Error
	if err != nil {
		return Assignment{}, logs.ErrWrapCtx(ctx, err, "find case files")
	}

	if len(files) != 2 {
		return Assignment{}, errors.New("invalid case files count")
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

	return toAssignment(assignment, course, user, caseInputFile, caseOutputFile), err
}

type FindAllAssignmentsRequest struct {
	core.Pagination
	CourseID ulids.ULID
}

type FindAllAssignmentsResponse struct {
	Assignments []Assignment
	core.Pagination
}

func (AssignmentReader) FindAll(ctx context.Context, tx *gorm.DB, req FindAllAssignmentsRequest) (FindAllAssignmentsResponse, error) {
	sqldb, err := dbconn.DBTxFromGorm(tx)
	if err != nil {
		return FindAllAssignmentsResponse{}, logs.ErrWrapCtx(ctx, err, "AssignmentReader: FindAll: cast dbtx")
	}

	query := xsqlc.New(sqldb)

	course, err := query.FindCourseByID(ctx, req.CourseID.String())
	if err != nil {
		return FindAllAssignmentsResponse{}, fmt.Errorf("find course: %w", err)
	}

	assignments, err := query.FindAllAssignmentsByCourseID(ctx, xsqlc.FindAllAssignmentsByCourseIDParams{
		CourseID:   req.CourseID.String(),
		PageOffset: req.Offset(),
		PageLimit:  req.Limit,
	})
	if err != nil {
		return FindAllAssignmentsResponse{}, logs.ErrWrapCtx(ctx, err, "AssignmentReader: FindAll: query")
	}

	count, err := query.CountAllAssignmentsByCourse(ctx, req.CourseID.String())
	if err != nil {
		return FindAllAssignmentsResponse{}, err
	}

	userIDs := []string{}
	for _, assignment := range assignments {
		userIDs = append(userIDs, assignment.AssignedBy)
	}

	users := []dbmodel.User{}
	err = tx.Table("users").Where("id IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return FindAllAssignmentsResponse{}, err
	}

	userMap := make(map[string]dbmodel.User, len(users))
	for _, user := range users {
		userMap[user.ID.String()] = user
	}

	fileIDs := []string{}
	for _, assignment := range assignments {
		fileIDs = append(fileIDs, assignment.CaseInputFileID, assignment.CaseOutputFileID)
	}

	files := []dbmodel.File{}
	err = tx.Table("files").Where("id IN (?)", fileIDs).Find(&files).Error
	if err != nil {
		return FindAllAssignmentsResponse{}, err
	}

	fileMap := make(map[string]dbmodel.File, len(files))
	for _, file := range files {
		fileMap[file.ID.String()] = file
	}

	result := FindAllAssignmentsResponse{
		Pagination: core.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: int32(count),
		},
		Assignments: make([]Assignment, len(assignments)),
	}

	for i, assignment := range assignments {
		user := userMap[assignment.AssignedBy]
		fileInput := fileMap[assignment.CaseInputFileID]
		fileOutput := fileMap[assignment.CaseOutputFileID]

		asg := toAssignment(xsqlc.FindAssignmentByIDRow(assignment), course, user, fileInput, fileOutput)
		result.Assignments[i] = asg
	}

	return result, nil
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
	model xsqlc.FindAssignmentByIDRow,
	course xsqlc.Course,
	user dbmodel.User,
	inputFile dbmodel.File,
	outputFile dbmodel.File,
) Assignment {
	courseID, err := ulids.Parse(course.ID)
	if err != nil {
		panic(err)
	}

	return Assignment{
		ID:             uuid.MustParse(model.ID),
		CourseID:       courseID,
		Name:           model.Name,
		Description:    model.Description,
		Assigner:       toAssigner(user),
		CaseInputFile:  toCaseFile(inputFile),
		CaseOutputFile: toCaseFile(outputFile),
		TimestampMetadata: core.TimestampMetadata{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: null.Time{NullTime: model.DeletedAt},
		},
		Template:   model.Template,
		DeadlineAt: model.DeadlineAt,
		Course: Course{
			ID:          courseID,
			Name:        course.Name,
			Description: course.Description,
			IsActive:    course.IsActive,
		},
	}
}
