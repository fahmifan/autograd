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
	fileReader := assignments.FileReader{}

	now := time.Now()
	deadlineAt, err := time.Parse(time.RFC3339, req.Msg.GetDeadlineAt())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	assignment := assignments.Assignment{}
	caseStdinFileID, err := uuid.Parse(req.Msg.GetCaseInputFileId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	caseStdoutFileID, err := uuid.Parse(req.Msg.GetCaseOutputFileId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) error {
		assigner, err := assignerReader.FindByID(ctx, cmd.GormDB, authUser.UserID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: FindByID")
			return core.ErrInternalServer
		}

		fileIDs := []uuid.UUID{caseStdinFileID, caseStdoutFileID}
		caseFiles, err := fileReader.FindCaseFiles(ctx, cmd.GormDB, fileIDs)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: FindCaseFiles")
			return core.ErrInternalServer
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

		assignment, err = assignments.CreateAssignment(assignments.CreateAssignmentRequest{
			NewID:          uuid.New(),
			Name:           req.Msg.GetName(),
			Description:    req.Msg.GetDescription(),
			Assigner:       assigner,
			Now:            now,
			CaseInputFile:  caseInputFile,
			CaseOutputFile: caseOutputFile,
			DeadlineAt:     deadlineAt,
		})
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignmentWriter.Create(ctx, tx, assignment)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: Save")
			return core.ErrInternalServer
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
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !authUser.Role.Can(auth.UpdateAssignment) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	assignmentID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
			return core.ErrInternalServer
		}

		assignment, err := assignmentReader.FindByID(ctx, cmd.GormDB, assignmentID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: FindByID")
			return core.ErrInternalServer
		}

		fileIDs := []uuid.UUID{assignment.CaseInputFile.ID, assignment.CaseOutputFile.ID}
		caseFiles, err := fileReader.FindCaseFiles(ctx, cmd.GormDB, fileIDs)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: FindCaseFiles")
			return core.ErrInternalServer
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

		err = assignmentWriter.Create(ctx, tx, assignment)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: UpdateAssignment: Save")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return core.ProtoEmptyResponse, nil
}

func (cmd *AssignmentCmd) DeleteAssignment(ctx context.Context, req *connect.Request[autogradv1.DeleteByIDRequest]) (*connect.Response[autogradv1.Empty], error) {

	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !authUser.Role.Can(auth.DeleteAssignment) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	assignmentID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	assignmentReader := assignments.AssignmentReader{}
	assignmentWriter := assignments.AssignmentWriter{}
	now := time.Now()

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) error {
		assignment, err := assignmentReader.FindByID(ctx, cmd.GormDB, assignmentID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: DeleteAssignment: FindByID")
			return core.ErrInternalServer
		}

		assignment, err = assignment.Delete(now)
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignmentWriter.Create(ctx, tx, assignment)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateAssignment: Save")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return core.ProtoEmptyResponse, nil
}

func (cmd *AssignmentCmd) CreateSubmission(ctx context.Context, req *connect.Request[autogradv1.CreateSubmissionRequest]) (*connect.Response[autogradv1.CreatedResponse], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	isCreateForOther := req.Msg.GetSubmitterId() != authUser.UserID.String()
	isAllowCreateForOther := authUser.Role.Can(auth.CreateSubmissionForOther)
	if isCreateForOther && !isAllowCreateForOther {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	if !authUser.Role.Can(auth.CreateSubmission) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	assignmentID, err := uuid.Parse(req.Msg.GetAssignmentId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	submitterID, err := uuid.Parse(req.Msg.GetSubmitterId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	submissionFileID, err := uuid.Parse(req.Msg.GetSourceFileId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	now := time.Now()
	submission := assignments.Submission{}

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) (err error) {
		assignment, err := assignments.AssignmentReader{}.FindByID(ctx, cmd.GormDB, assignmentID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: AssignmentReader{}.FindByID")
			return core.ErrInternalServer
		}

		submitter, err := assignments.SubmitterReader{}.FindByID(ctx, cmd.GormDB, submitterID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmitterReader{}.FindByID")
			return core.ErrInternalServer
		}

		submissionFile, err := assignments.SubmissionFileReader{}.FindByID(ctx, cmd.GormDB, submissionFileID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmissionFileReader{}.FindByID")
			return core.ErrInternalServer
		}

		submission, err = assignments.CreateSubmission(assignments.CreateSubmissionRequest{
			NewID:          uuid.New(),
			Now:            now,
			Assignment:     assignment,
			Submitter:      submitter,
			SubmissionFile: submissionFile,
		})
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignments.SubmissionWriter{}.SaveNew(ctx, cmd.GormDB, &submission)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmissionWriter{}.Save")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	res := &connect.Response[autogradv1.CreatedResponse]{
		Msg: &autogradv1.CreatedResponse{
			Id:      submission.ID.String(),
			Message: "submission created",
		},
	}

	return res, nil
}

func (cmd *AssignmentCmd) UpdateSubmission(ctx context.Context, req *connect.Request[autogradv1.UpdateSubmissionRequest]) (*connect.Response[autogradv1.Empty], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	isUpdateForOther := req.Msg.GetSubmitterId() != authUser.UserID.String()
	isAllowCreateForOther := authUser.Role.Can(auth.CreateSubmissionForOther)
	if isUpdateForOther && !isAllowCreateForOther {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	if !authUser.Role.Can(auth.CreateSubmission) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	submitterID, err := uuid.Parse(req.Msg.GetSubmitterId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	submissionID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	submissionFileID, err := uuid.Parse(req.Msg.GetSourceFileId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	now := time.Now()
	submission := assignments.Submission{}

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) (err error) {
		submission, err = assignments.SubmissionReader{}.FindByID(ctx, cmd.GormDB, submissionID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: AssignmentReader{}.FindByID")
			return core.ErrInternalServer
		}

		submitter, err := assignments.SubmitterReader{}.FindByID(ctx, cmd.GormDB, submitterID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmitterReader{}.FindByID")
			return core.ErrInternalServer
		}

		submissionFile, err := assignments.SubmissionFileReader{}.FindByID(ctx, cmd.GormDB, submissionFileID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmissionFileReader{}.FindByID")
			return core.ErrInternalServer
		}

		submission, err = submission.Update(assignments.UpdateSubmissionRequest{
			Now:            now,
			SubmissionFile: submissionFile,
			Submitter:      submitter,
		})
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignments.SubmissionWriter{}.Save(ctx, cmd.GormDB, &submission)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmissionWriter{}.Save")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return core.ProtoEmptyResponse, nil
}

func (cmd *AssignmentCmd) DeleteSubmission(ctx context.Context, req *connect.Request[autogradv1.DeleteByIDRequest]) (*connect.Response[autogradv1.Empty], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !authUser.Role.Can(auth.DeleteSubmission) {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	submissionID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	now := time.Now()
	submission := assignments.Submission{}

	err = core.Transaction(cmd.Ctx, func(tx *gorm.DB) (err error) {
		submission, err = assignments.SubmissionReader{}.FindByID(ctx, cmd.GormDB, submissionID)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: AssignmentReader{}.FindByID")
			return core.ErrInternalServer
		}

		isDeleteForOther := !submission.IsOwner(authUser.UserID)
		isAllowDeleteForOther := authUser.Role.Can(auth.DeleteSubmissionForOther)
		if isDeleteForOther && !isAllowDeleteForOther {
			return connect.NewError(connect.CodePermissionDenied, nil)
		}

		submission, err = submission.Delete(now)
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = assignments.SubmissionWriter{}.Delete(ctx, cmd.GormDB, &submission)
		if err != nil {
			logs.ErrCtx(ctx, err, "AssignmentCmd: CreateSubmission: SubmissionWriter{}.Save")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return core.ProtoEmptyResponse, nil
}
