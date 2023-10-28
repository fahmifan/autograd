package assignments_query

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/assignments"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/google/uuid"
)

type AssignmentsQuery struct {
	*core.Ctx
}

func (query *AssignmentsQuery) FindAllAssignments(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllAssignmentsRequest],
) (*connect.Response[autogradv1.FindAllAssignmentsResponse], error) {
	res, err := assignments.AssignmentReader{}.FindAll(ctx, query.GormDB, assignments.FindAllAssignmentsRequest{
		Page:  req.Msg.GetPage(),
		Limit: req.Msg.GetLimit(),
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "AssignmentsQuery: FindAllAssignments: FindAll")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[autogradv1.FindAllAssignmentsResponse]{
		Msg: &autogradv1.FindAllAssignmentsResponse{
			Assignments: toAssignmentProtos(res.Assignments),
			PaginationMetadata: &autogradv1.PaginationMetadata{
				Total: res.Count,
				Page:  req.Msg.GetPage(),
				Limit: req.Msg.GetLimit(),
			},
		},
	}, nil
}

func (query *AssignmentsQuery) FindAssignment(
	ctx context.Context,
	req *connect.Request[autogradv1.FindByIDRequest],
) (*connect.Response[autogradv1.Assignment], error) {
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

	return &connect.Response[autogradv1.Submission]{
		Msg: toSubmissionProto(submission),
	}, nil
}

func toSubmissionProto(submission assignments.Submission) *autogradv1.Submission {
	return &autogradv1.Submission{
		Id:         submission.ID.String(),
		Assignment: toAssignmentProto(submission.Assignment),
		Submitter:  toSubmitterProto(submission.Submitter),
		SubmissionFile: &autogradv1.SubmissionFile{
			Id:  submission.SourceFile.ID.String(),
			Url: submission.SourceFile.URL,
		},
		TimestampMetadata: submission.ProtoTimestampMetadata(),
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
		TimestampMetadata: assignment.ProtoTimestampMetadata(),
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
