package student_assignment_query

import (
	"context"
	"errors"
	"io"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/mediastore/mediastore_query"
	"github.com/fahmifan/autograd/pkg/core/student_assignment"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/google/uuid"
)

type StudentAssignmentQuery struct {
	*core.Ctx
}

func (query *StudentAssignmentQuery) FindAllStudentAssignments(ctx context.Context, req *connect.Request[autogradv1.FindAllStudentAssignmentsRequest]) (
	*connect.Response[autogradv1.FindAllStudentAssignmentsResponse], error,
) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewAssignment) {
		return nil, core.ErrPermissionDenied
	}

	var (
		from time.Time
		to   time.Time
		err  error
	)

	if req.Msg.GetFromDate() != "" {
		from, err = time.Parse(time.RFC3339, req.Msg.GetFromDate())
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid from date"))
		}
	}

	if req.Msg.GetToDate() != "" {
		to, err = time.Parse(time.RFC3339, req.Msg.GetToDate())
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid to date"))
		}
	}

	reader := student_assignment.StudentAssignmentReader{}
	res, err := reader.FindAllAssignments(ctx, query.GormDB, student_assignment.FindAllAssignmentRequest{
		PaginationRequest: core.PaginationRequestFromProto(req.Msg.GetPaginationRequest()),
		From:              from,
		To:                to,
		StudentID:         authUser.UserID,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "StudentAssignmentQuery: FindAllStudentAssignments: FindAllAssignments")
		return nil, core.ErrInternalServer
	}

	return &connect.Response[autogradv1.FindAllStudentAssignmentsResponse]{
		Msg: &autogradv1.FindAllStudentAssignmentsResponse{
			Assignments:        toStudentAssignmentProtos(res.Assignments),
			PaginationMetadata: res.ProtoPagination(),
		},
	}, nil
}

func (query *StudentAssignmentQuery) FindStudentAssignment(ctx context.Context, req *connect.Request[autogradv1.FindByIDRequest]) (
	*connect.Response[autogradv1.StudentAssignment], error,
) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewAssignment) {
		return nil, core.ErrPermissionDenied
	}

	id, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid id"))
	}

	reader := student_assignment.StudentAssignmentReader{}
	res, err := reader.FindByID(ctx, query.GormDB, student_assignment.FindStudentAssignmentByIDRequest{
		ID:            id,
		StudentID:     authUser.UserID,
		WithStudentID: true,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "StudentAssignmentQuery: FindStudentAssignment: FindByID")
		return nil, core.ErrInternalServer
	}

	var submissionBuf []byte
	if res.HasSubmission {
		mediaStoreQuery := mediastore_query.MediaStoreQuery{Ctx: query.Ctx}
		submissionMedia, err := mediaStoreQuery.InternalFindMediaFile(ctx, mediastore_query.InternalFindMediaFileRequest{
			ID: res.Submission.SubmissionFileID,
		})
		if err != nil {
			logs.ErrCtx(ctx, err, "StudentAssignmentQuery: FindStudentAssignment: InternalFindMediaFile")
			return nil, core.ErrInternalServer
		}
		defer submissionMedia.BodyCloser.Close()

		submissionBuf, err = io.ReadAll(submissionMedia.BodyCloser)
		if err != nil {
			logs.ErrCtx(ctx, err, "StudentAssignmentQuery: FindStudentAssignment: io.ReadAll")
			return nil, core.ErrInternalServer
		}
	}

	return &connect.Response[autogradv1.StudentAssignment]{
		Msg: toStudentAssignmentProto(res, submissionBuf),
	}, nil
}

func toStudentAssignmentProtos(assignments []student_assignment.StudentAssignment) []*autogradv1.StudentAssignment {
	var assignmentProtos []*autogradv1.StudentAssignment
	for _, assignment := range assignments {
		assignmentProtos = append(assignmentProtos, toStudentAssignmentProto(assignment, nil))
	}
	return assignmentProtos
}

func toStudentAssignmentProto(assignment student_assignment.StudentAssignment, submissionCode []byte) *autogradv1.StudentAssignment {
	return &autogradv1.StudentAssignment{
		Id:           assignment.ID.String(),
		Name:         assignment.Name,
		Description:  assignment.Description,
		CodeTemplate: assignment.CodeTemplate,
		AssignerId:   assignment.Assigner.ID.String(),
		AssignerName: assignment.Assigner.Name,
		UpdatedAt:    assignment.UpdatedAt.Format(time.RFC3339),
		DeadlineAt:   assignment.DeadlineAt.Format(time.RFC3339),
		Submission: &autogradv1.StudentAssignment_Submission{
			Id:             assignment.Submission.ID.String(),
			SubmissionCode: string(submissionCode),
			Grade:          int32(assignment.Submission.Grade),
			UpdatedAt:      assignment.Submission.UpdatedAt.Format(time.RFC3339),
			IsGraded:       assignment.Submission.IsGraded,
		},
		HasSubmission: assignment.HasSubmission,
	}
}
