package assignments_query

import (
	"context"
	"io"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/assignments"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/mediastore/mediastore_query"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/ulids"
	"github.com/google/uuid"
)

type AssignmentsQuery struct {
	*core.Ctx
}

func (query *AssignmentsQuery) FindAllAssignments(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllAssignmentsRequest],
) (*connect.Response[autogradv1.FindAllAssignmentsResponse], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewAnyAssignments) {
		return nil, core.ErrPermissionDenied
	}

	courseID, err := ulids.Parse(req.Msg.GetCourseId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res, err := assignments.AssignmentReader{}.FindAll(ctx, query.GormDB, assignments.FindAllAssignmentsRequest{
		Pagination: core.PaginationRequestFromProto(req.Msg.GetPaginationRequest()),
		CourseID:   courseID,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "AssignmentsQuery: FindAllAssignments: FindAll")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[autogradv1.FindAllAssignmentsResponse]{
		Msg: &autogradv1.FindAllAssignmentsResponse{
			Assignments:        toAssignmentProtos(res.Assignments),
			PaginationMetadata: res.ProtoPagination(),
		},
	}, nil
}

func (query *AssignmentsQuery) FindAssignment(
	ctx context.Context,
	req *connect.Request[autogradv1.FindByIDRequest],
) (*connect.Response[autogradv1.Assignment], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewAnyAssignments) {
		return nil, core.ErrPermissionDenied
	}

	assignmentID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	assignmentReader := assignments.AssignmentReader{}
	assignment, err := assignmentReader.FindByID(ctx, query.GormDB, assignmentID)
	if core.IsDBNotFoundErr(err) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		logs.ErrCtx(ctx, err, "AssignmentsQuery: FindAssignment: FindByID")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[autogradv1.Assignment]{
		Msg: toAssignmentProto(assignment),
	}, nil
}

func (query *AssignmentsQuery) FindSubmission(
	ctx context.Context,
	req *connect.Request[autogradv1.FindByIDRequest],
) (*connect.Response[autogradv1.Submission], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewAnySubmissions) {
		return nil, core.ErrPermissionDenied
	}

	submissionID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	submReader := assignments.SubmissionReader{}
	submission, err := submReader.FindByID(ctx, query.GormDB, submissionID)
	if core.IsDBNotFoundErr(err) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		logs.ErrCtx(ctx, err, "AssignmentsQuery: FindSubmission: FindByID")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var submissionBuf []byte
	if submission.SourceFile.ID != uuid.Nil {
		mediaStoreQuery := mediastore_query.MediaStoreQuery{Ctx: query.Ctx}
		submissionMedia, err := mediaStoreQuery.InternalFindMediaFile(ctx, mediastore_query.InternalFindMediaFileRequest{
			ID: submission.SourceFile.ID,
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

	return &connect.Response[autogradv1.Submission]{
		Msg: toSubmissionProto(submission, submissionBuf),
	}, nil
}

func (query *AssignmentsQuery) FindAllSubmissionForAssignment(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllSubmissionsForAssignmentRequest],
) (*connect.Response[autogradv1.FindAllSubmissionsForAssignmentResponse], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewAnySubmissions) {
		return nil, core.ErrPermissionDenied
	}

	assignment := dbmodel.Assignment{}
	err := query.GormDB.Where("id = ?", req.Msg.GetAssignmentId()).Take(&assignment).Error
	if err != nil && core.IsDBNotFoundErr(err) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	assigner := dbmodel.User{}
	err = query.GormDB.Where("id = ?", assignment.AssignedBy).Take(&assigner).Error
	if err != nil && core.IsDBNotFoundErr(err) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	var submissions []dbmodel.Submission
	err = query.GormDB.
		Model(&dbmodel.Submission{}).
		Select("id", "submitted_by").
		Where("assignment_id = ?", assignment.ID).Find(&submissions).Error
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	submitterIDs := make([]uuid.UUID, len(submissions))
	for i, submission := range submissions {
		submitterIDs[i] = submission.SubmittedBy
	}

	var submitters []dbmodel.User
	err = query.GormDB.
		Select("id", "name").
		Where("id IN ?", submitterIDs).Find(&submitters).Error
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	submitterMap := make(map[uuid.UUID]string)
	for _, submitter := range submitters {
		submitterMap[submitter.ID] = submitter.Name
	}

	submissionRes := make([]*autogradv1.FindAllSubmissionsForAssignmentResponse_Submission, len(submissions))
	for i := range submissions {
		submissionRes[i] = &autogradv1.FindAllSubmissionsForAssignmentResponse_Submission{
			Id:            submissions[i].ID.String(),
			SubmitterId:   submissions[i].SubmittedBy.String(),
			SubmitterName: submitterMap[submissions[i].SubmittedBy],
		}
	}

	res := &connect.Response[autogradv1.FindAllSubmissionsForAssignmentResponse]{
		Msg: &autogradv1.FindAllSubmissionsForAssignmentResponse{
			Submissions:    submissionRes,
			AssignerId:     assigner.ID.String(),
			AssignerName:   assigner.Name,
			AssignmentId:   assignment.ID.String(),
			AssignmentName: assignment.Name,
		},
	}

	return res, nil
}

func toSubmissionProto(submission assignments.Submission, submissionBuf []byte) *autogradv1.Submission {
	return &autogradv1.Submission{
		Id:         submission.ID.String(),
		Assignment: toAssignmentProto(submission.Assignment),
		Submitter:  toSubmitterProto(submission.Submitter),
		SubmissionFile: &autogradv1.SubmissionFile{
			Id:  submission.SourceFile.ID.String(),
			Url: submission.SourceFile.URL,
		},
		TimestampMetadata: submission.ProtoTimestampMetadata(),
		SubmissionCode:    string(submissionBuf),
	}
}

func toSubmitterProto(submitter assignments.Submitter) *autogradv1.Submitter {
	return &autogradv1.Submitter{
		Id:   submitter.ID.String(),
		Name: submitter.Name,
	}
}

func toAssignmentProtos(assignments []assignments.Assignment) []*autogradv1.Assignment {
	var result []*autogradv1.Assignment
	for _, assignment := range assignments {
		result = append(result, toAssignmentProto(assignment))
	}
	return result
}

func toAssignmentProto(assignment assignments.Assignment) *autogradv1.Assignment {
	return &autogradv1.Assignment{
		Id:                assignment.ID.String(),
		Name:              assignment.Name,
		Description:       assignment.Description,
		Template:          assignment.Template,
		TimestampMetadata: assignment.ProtoTimestampMetadata(),
		DeadlineAt:        assignment.DeadlineAt.Format(time.RFC3339),
		Assigner: &autogradv1.Assigner{
			Id:   assignment.Assigner.ID.String(),
			Name: assignment.Assigner.Name,
		},
		CaseInputFile: &autogradv1.AssignmentFile{
			Id:                assignment.CaseInputFile.ID.String(),
			Url:               assignment.CaseInputFile.URL,
			TimestampMetadata: assignment.CaseInputFile.ProtoTimestampMetadata(),
		},
		CaseOutputFile: &autogradv1.AssignmentFile{
			Id:                assignment.CaseOutputFile.ID.String(),
			Url:               assignment.CaseOutputFile.URL,
			TimestampMetadata: assignment.CaseInputFile.ProtoTimestampMetadata(),
		},
	}
}
