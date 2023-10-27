package assignments_cmd

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/assignments"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type AssignmentCmd struct {
	*core.Ctx
}

func (cmd *AssignmentCmd) CreateAssignment(ctx context.Context, req *connect.Request[autogradv1.CreateAssignmentRequest]) (*connect.Response[autogradv1.CreatedResponse], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !authUser.Role.Can(auth.CreateAssignment) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	assignerReader := assignments.AssignerReader{}
	assignmentWriter := assignments.AssignmentWriter{}

	assignment := assignments.Assignment{}
	err := core.Transaction(cmd.Ctx, func(tx *gorm.DB) error {
		assigner, err := assignerReader.FindByID(ctx, cmd.GormDB, authUser.UserID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: FindByID")
			return connect.NewError(connect.CodeInternal, err)
		}

		assignment, err = assignments.CreateAssignment(assignments.CreateAssignmentRequest{
			NewID:       uuid.New(),
			Name:        req.Msg.GetName(),
			Description: req.Msg.GetDescription(),
			Assigner:    assigner,
		})
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignmentWriter.Save(ctx, tx, assignment)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: Save")
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &connect.Response[autogradv1.CreatedResponse]{
		Msg: &autogradv1.CreatedResponse{
			Id:      assignment.ID.String(),
			Message: "assignment created",
		},
	}, nil
}

func (cmd *AssignmentCmd) UpdateAssignment(ctx context.Context, req *connect.Request[autogradv1.UpdateAssignmentRequest]) (*connect.Response[autogradv1.Empty], error) {
	assignmentID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !authUser.Role.Can(auth.UpdateAssignment) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	now := time.Now()

	assignerReader := assignments.AssignerReader{}
	assignmentReader := assignments.AssignmentReader{}
	assignmentWriter := assignments.AssignmentWriter{}
	fileReader := assignments.FileReader{}

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) error {
		assigner, err := assignerReader.FindByID(ctx, cmd.GormDB, authUser.UserID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: FindByID")
			return connect.NewError(connect.CodeInternal, err)
		}

		assignment, err := assignmentReader.FindByID(ctx, cmd.GormDB, assignmentID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: FindByID")
			return connect.NewError(connect.CodeInternal, err)
		}

		fileIDs := []uuid.UUID{assignment.CaseInputFile.ID, assignment.CaseOutputFile.ID}
		caseFiles, err := fileReader.FindCaseFiles(ctx, cmd.GormDB, fileIDs)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: FindCaseFiles")
			return connect.NewError(connect.CodeInternal, err)
		}

		caseInputFile, _, found := lo.FindIndexOf(caseFiles, func(file assignments.CaseFile) bool {
			return file.Type == dbmodel.FileTypeAssignmentCaseInput
		})
		if !found {
			return errors.New("case input file not found")
		}

		caseOutputFile, _, found := lo.FindIndexOf(caseFiles, func(file assignments.CaseFile) bool {
			return file.Type == dbmodel.FileTypeAssignmentCaseOutput
		})
		if !found {
			return errors.New("case output file not found")
		}

		assignment, err = assignment.Update(assignments.UpdateAssignmentRequest{
			Now:            now,
			Name:           req.Msg.GetName(),
			Description:    req.Msg.GetDescription(),
			Assigner:       assigner,
			CaseInputFile:  caseInputFile,
			CaseOutputFile: caseOutputFile,
		})
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignmentWriter.Save(ctx, tx, assignment)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: Save")
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return core.EmptyResponse, nil
}

func (cmd *AssignmentCmd) DeleteAssignment(ctx context.Context, req *connect.Request[autogradv1.DeleteByIDRequest]) (*connect.Response[autogradv1.Empty], error) {
	assignmentID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !authUser.Role.Can(auth.CreateAssignment) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	assignmentReader := assignments.AssignmentReader{}
	assignmentWriter := assignments.AssignmentWriter{}
	now := time.Now()

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) error {
		assignment, err := assignmentReader.FindByID(ctx, cmd.GormDB, assignmentID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: DeleteAssignment: FindByID")
			return connect.NewError(connect.CodeInternal, err)
		}

		assignment, err = assignment.Delete(now)
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignmentWriter.Save(ctx, tx, assignment)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: Save")
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return core.EmptyResponse, nil
}
